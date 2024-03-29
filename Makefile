.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_user uc-users/deliveries/create-user/create_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sign_in uc-users/deliveries/sign-in/sign_in.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/confirm_user uc-users/deliveries/confirm/confirm_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_user uc-users/deliveries/get-user/get_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/remove_notification uc-users/deliveries/remove-notification/remove-notification.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/send_friend_req uc-users/deliveries/send-friend-req/send_friend_req.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/accept_friend_req uc-users/deliveries/accept-friend-req/accept_friend_req.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/unfriend uc-users/deliveries/unfriend/unfriend.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/search_friends uc-users/deliveries/search-friends/search_friends.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/get_recent_swing_uploads uc-uploads/deliveries/get_recent_swing_uploads/get_recent_swing_uploads.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_swing_upload uc-uploads/deliveries/create_swing_upload/create_swing_upload.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_upload_clips uc-uploads/deliveries/update_upload_clips/update_upload_clips.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_upload_swings uc-uploads/deliveries/update_upload_swings/update_upload_swings.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/get_user_albums uc-albums/deliveries/get_user_albums/get_user_albums.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_album uc-albums/deliveries/get_album/get_album.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_album uc-albums/deliveries/create_album/create_album.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete_album uc-albums/deliveries/delete_album/delete_album.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_album uc-albums/deliveries/update_album/update_album.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post_comment uc-albums/deliveries/post_comment/post_comment.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/create_comment_flag uc-moderation/deliveries/create_comment_flag/create_comment_flag.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_album_flag uc-moderation/deliveries/create_album_flag/create_album_flag.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_comment_flag uc-moderation/deliveries/update_comment_flag/update_comment_flag.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_album_flag uc-moderation/deliveries/update_album_flag/update_album_flag.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/recent_flagged_comments uc-moderation/deliveries/recent_flagged_comments/recent_flagged_comments.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/recent_flagged_albums uc-moderation/deliveries/recent_flagged_albums/recent_flagged_albums.go

	cp ./id_rsa bin/id_rsa
	cp ./id_rsa.pub bin/id_rsa.pub

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
