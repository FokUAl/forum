build:
	docker image build -f Dockerfile -t forum_image .
run:
	docker container run -p 4888:4888 --detach --name application forum_image