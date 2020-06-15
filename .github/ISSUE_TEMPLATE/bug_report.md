---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

version of yq: 
operating system:

**Input Yaml**
Concise yaml document(s) (as simple as possible to show the bug)
data1.yml:
```yaml
this: should really work
```

data2.yml:
```yaml
but: it strangely didn't
```

**Command**
The command you ran:
```
yq merge data1.yml data2.yml
```

**Actual behavior**

```yaml
cat: meow
```

**Expected behavior**

```yaml
this: should really work
but: it strangely didn't
```

**Additional context**
Add any other context about the problem here.
