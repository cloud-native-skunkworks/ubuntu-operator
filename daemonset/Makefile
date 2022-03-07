all:
	docker build -t "tibbar/ubuntu-machine-controller:latest" .
	docker push "tibbar/ubuntu-machine-controller:latest"
test:
	go run main.go --socketPath /home/alex/Code/uo2.sock --module=nvme_core=foo --module=rfcomm --apt="build-essentials" --snap=microK8s=classic --kubeconfig=/home/alex/Code/config