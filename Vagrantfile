# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  if Vagrant.has_plugin?('landrush')
    config.landrush.enabled = true
    config.landrush.tld = 'dev'
    config.landrush.guest_redirect_dns = false
  end

  config.vm.define 'development' do |base|
    base.vm.provider 'virtualbox' do |v|
      v.memory = 1024
      v.cpus = 2
    end
    base.vm.box = 'ubuntu/trusty64'
    base.vm.hostname = 'development.vagrant.dev'
    base.vm.network 'private_network', type: 'dhcp'

    base.vm.network "forwarded_port", guest: 8080, host: 8080

    # base.vm.synced_folder "~/Documents/workspace", "/workspace"
  end

  config.vm.provision "shell", inline: "curl -sSL https://get.docker.com/ubuntu/ | sudo sh"
  config.vm.provision "shell", inline: "sudo gpasswd -a vagrant docker"
  config.vm.provision "shell", inline: "sudo service docker restart"
  # config.vm.provision "shell", inline: "apt-get install -y unzip openjdk-6-jre-headless"
  config.vm.provision "shell", inline: "curl -L -o packer.zip https://dl.bintray.com/mitchellh/packer/packer_0.8.1_linux_amd64.zip"
  config.vm.provision "shell", inline: "unzip packer.zip -d packer"
  config.vm.provision "shell", inline: "mv -f packer /usr/local/packer"
  config.vm.provision "shell", inline: "rm packer.zip"
  config.vm.provision "shell", inline: "sudo su -c \"echo \\\"export PATH=$PATH:/usr/local/packer\\\" > /etc/profile.d/packer.sh\""

end
