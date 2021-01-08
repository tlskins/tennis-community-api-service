.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_user uc-users/deliveries/create-user/create_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sign_in uc-users/deliveries/sign-in/sign_in.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/confirm_user uc-users/deliveries/confirm/confirm_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_user uc-users/deliveries/get-user/get_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/clear_notifications uc-users/deliveries/clear-notifications/clear_notifications.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/get_recent_swing_uploads uc-uploads/deliveries/get_recent_swing_uploads/get_recent_swing_uploads.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_swing_upload uc-uploads/deliveries/create_swing_upload/create_swing_upload.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_upload_clips uc-uploads/deliveries/update_upload_clips/update_upload_clips.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_upload_swings uc-uploads/deliveries/update_upload_swings/update_upload_swings.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/get_user_albums uc-albums/deliveries/get_user_albums/get_user_albums.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_album uc-albums/deliveries/get_album/get_album.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_album uc-albums/deliveries/update_album/update_album.go

	cp ./id_rsa bin/id_rsa
	cp ./id_rsa.pub bin/id_rsa.pub

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
