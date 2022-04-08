package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/xanzy/go-gitlab"
)

// Checks for the latest yq version and if there is a new one,
// creates a new branch on alpine gitlab with the changes to update to the new version.
// Make sure to already have a fork of alpine/aports ready and an access token (api) created.
func main() {
	// Injected (Env) variables
	gitlabToken := os.Getenv("ALPINE_GITLAB_ACCESS_TOKEN")
	gitlabForkRepo := os.Getenv("ALPINE_GITLAB_FORK_REPO")
	gitUsername := os.Getenv("GIT_USERNAME")
	gitEmail := os.Getenv("GIT_EMAIL")
	autocreateMergeRequest, _ := strconv.ParseBool(os.Getenv("AUTOCREATE_MERGE_REQUEST"))

	// Debug
	//gitlabToken = "xxx"
	//gitlabForkRepo = "xxx/aports"
	//gitUsername = "xxx xxx"
	//gitEmail = "xxx@xxx.xxx"
	//autocreateMergeRequest = true

	// Variables
	gitlabBaseUrl := "gitlab.alpinelinux.org"
	gitlabRemoteRepo := "alpine/aports"
	gitlabMainBranch := "master"
	tempCheckoutPath := "./temp_checkout"

	// Cleanup
	err := os.RemoveAll(tempCheckoutPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Checkout a fresh copy of the repo
	fmt.Println("Cloning and updating repo...")
	cloneUrl := fmt.Sprintf("https://%s/%s.git", gitlabBaseUrl, gitlabForkRepo)
	err = exec.Command("git", "clone", cloneUrl, tempCheckoutPath).Run() // #nosec G204
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Switch to work in this directory from now on
	err = os.Chdir(tempCheckoutPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Register the upstream
	err = exec.Command("git", "remote", "add", "upstream", fmt.Sprintf("https://%s/%s.git", gitlabBaseUrl, gitlabRemoteRepo)).Run() // #nosec G204
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Update the forks main branch
	err = exec.Command("git", "checkout", gitlabMainBranch).Run() // #nosec G204
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = exec.Command("git", "fetch", "upstream").Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = exec.Command("git", "merge", fmt.Sprintf("upstream/%s", gitlabMainBranch)).Run() // #nosec G204
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Get information about the latest version and its hash
	fmt.Println("Checking for new version...")
	packageUser := "mikefarah"
	packageRepo := "yq"
	latestVersion := getLatestReleaseFromGithub(packageUser, packageRepo)
	hash := calculateMd512FromUrl(fmt.Sprintf("https://github.com/%s/%s/archive/v%s.tar.gz", packageUser, packageRepo, latestVersion))
	fmt.Printf("Found version %s\n", latestVersion)

	// Prepare the new branch
	branchName := fmt.Sprintf("feature/%s_%s_update", packageRepo, latestVersion)
	err = exec.Command("git", "checkout", "-b", branchName).Run() // #nosec G204
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Modify the file
	fmt.Println("Updating file...")
	apkFilePath := "./community/yq/APKBUILD"
	fileContent, err := os.ReadFile(apkFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	regVersion, _ := regexp.Compile("pkgver=(.*)")
	regChecksum, _ := regexp.Compile("(?m)sha512sums=\"(?:.|\n)*?\"")
	newContent := regVersion.ReplaceAll(fileContent, []byte(fmt.Sprintf("pkgver=%s", latestVersion)))
	newContent = regChecksum.ReplaceAll(newContent, []byte(fmt.Sprintf("sha512sums=\"\n%s  yq-%s.tar.gz\n\"", hash, latestVersion)))
	err = os.WriteFile(apkFilePath, []byte(newContent), 0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Check for changes and push the new branch
	out, err := exec.Command("git", "ls-files", "--modified").Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	updateVersionText := fmt.Sprintf("Updated yq to %s", latestVersion)
	if len(out) > 0 {
		fmt.Println("Changes detected, pushing new branch")
		err = exec.Command("git", "add", ".").Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = exec.Command("git", "-c", fmt.Sprintf("user.name='%s'", gitUsername), "-c", fmt.Sprintf("user.email='%s'", gitEmail), "commit", "-m", updateVersionText).Run() // #nosec G204
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = exec.Command("git", "push", "-u", "--force", fmt.Sprintf("https://gitlab-ci-token:%s@%s/%s.git", gitlabToken, gitlabBaseUrl, gitlabForkRepo), branchName).Run() // #nosec G204
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Pushed new branch. Branch url: https://%s/%s/-/tree/%s\n", gitlabBaseUrl, gitlabForkRepo, branchName)
	} else {
		fmt.Println("No changes detected, nothing to do.")
	}

	// Create a merge request for the new branch
	if autocreateMergeRequest {
		fmt.Println("Creating merge request")
		git, _ := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(fmt.Sprintf("https://%s/api/v4", gitlabBaseUrl)))
		projectRemote, _, _ := git.Projects.GetProject(gitlabRemoteRepo, &gitlab.GetProjectOptions{})
		projectFork, _, _ := git.Projects.GetProject(gitlabForkRepo, &gitlab.GetProjectOptions{})
		targetProjectId := projectRemote.ID
		_, _, err := git.MergeRequests.CreateMergeRequest(projectFork.ID, &gitlab.CreateMergeRequestOptions{
			Title:           &updateVersionText,
			SourceBranch:    &branchName,
			TargetBranch:    &gitlabMainBranch,
			TargetProjectID: &targetProjectId,
		})
		if err != nil {
			fmt.Println("Error creating Merge Request.")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func getLatestReleaseFromGithub(user string, repo string) string {
	var githubApiUrl = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", user, repo)

	var jsonResponse, _ = http.Get(githubApiUrl) // #nosec G107
	defer jsonResponse.Body.Close()

	var jsonData, _ = ioutil.ReadAll(jsonResponse.Body)
	var versionString, _ = jsonparser.GetString(jsonData, "tag_name")

	return versionString[1:]
}

func calculateMd512FromUrl(url string) string {
	var response, _ = http.Get(url) // #nosec G107
	defer response.Body.Close()

	var data, _ = ioutil.ReadAll(response.Body)
	var hash = sha512.Sum512(data)

	return hex.EncodeToString(hash[:])
}
