all:
	docker build -t "tibbar/ubuntu-machine-controller:latest" .
	docker push "tibbar/ubuntu-machine-controller:latest"
test:
	/home/linuxbrew/.linuxbrew/bin/go run main.go --socketPath /home/alex/Code/uo2.sock --module=nvme_core=foo --module=rfcomm --apt="build-essential" --snap=microk8s=classic --kubeconfig=/home/alex/Code/config
