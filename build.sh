operatingsystems=("linux" "windows" "darwin")
archs=("amd64" "arm64")

for os in ${operatingsystems[@]}; do
	for arch in ${archs[@]}; do
		bin_name="bin/$os-$arch-webserver"
		if [ "$os" = "windows" ]; then bin_name="$bin_name.exe"; fi
		GOOS=$os GOARCH=$arch go build -o $bin_name
		echo "Built $bin_name"
	done
done