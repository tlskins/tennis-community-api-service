.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	# env GOOS=linux go build -ldflags="-s -w" -o bin/upload_clips uc-uploads/deliveries/upload_clips.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_user uc-users/deliveries/create-user/create_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sign_in uc-users/deliveries/sign-in/sign_in.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/confirm_user uc-users/deliveries/confirm/confirm_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_swing_upload uc-uploads/deliveries/create_swing_upload/create_swing_upload.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_swing_upload_url uc-uploads/deliveries/get_swing_upload_url/get_swing_upload_url.go
	cp ./id_rsa bin/id_rsa
	cp ./id_rsa.pub bin/id_rsa.pub

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
