GREEN='\033[0;32m'
NC='\033[0m'

echo "${GREEN}---Old---${NC}"
yq $@ > /tmp/yq-old-output
cat /tmp/yq-old-output

echo "${GREEN}---New---${NC}"
./yq $@ > /tmp/yq-new-output
cat /tmp/yq-new-output

echo "${GREEN}---Diff---${NC}"
colordiff /tmp/yq-old-output /tmp/yq-new-output