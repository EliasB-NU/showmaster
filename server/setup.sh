# Update the system
sudo apt update -y && sudo apt upgrade -y

# Install Docker
sudo curl -fsSL https://get.docker.com -o get-docker.sh
sudo bash get-docker.sh

# Create Folder & get docker psql running
sudo mkdir /opt/stacks/psql && cd /opt/stacks/psql
sudo wget https://github.com/EliasB-NU/showmaster/blob/b4cf62a82f8fd229f43c71b65e6b334bc4332533/server/docker-compose.yml
docker compose up -d

# Create Folder for Project and clone it
sudo mkdir /opt/ && cd /opt/
sudo git clone https://github.com/EliasB-NU/showmaster.git

# Install Go (This setup is for ARM64 ONLY)
wget https://go.dev/dl/go1.22.3.linux-arm64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-arm64.tar.gz
sudo rm -rf go1.22.3.linux-arm64.tar.gz
sudo echo export PATH=$PATH:/usr/local/go/bin >> /etc/profile

# Compile Project
sudo /usr/local/go/bin/go build src/main.go

# Setup Systemctl
sudo mv /server/showmaster.service /etc/systemd/system/showmaster.service
sudo systemctl enable --now showmaster.service