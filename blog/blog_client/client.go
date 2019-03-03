package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/StanislawAbyszkin/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		AuthorId: "Stanislaw",
		Title:    "My third blog",
		Content:  "Content of my third blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error; %v\n", err)
	}
	fmt.Printf("Blog has been created: %v", res.GetBlog())
	blogID := res.GetBlog().GetId()

	// read blog
	fmt.Println("Reading the blog")

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5c7bb655946eca1c122f62cb"})
	if err != nil {
		log.Printf("Erro while reading; %v\n", err)
	}

	readRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		log.Printf("Erro while reading; %v\n", err)
	}
	log.Printf("Read blog with id: %v blog: %v\n", blogID, readRes.GetBlog())

	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed author",
		Title:    "My First Blog (Edited)",
		Content:  "Content of first Blog (Edited)",
	}
	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		log.Printf("Erro while updaing; %v\n", err)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete Blog
	fmt.Println("Deleting blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		log.Fatalf("Error while deleting; %v\n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list Blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			log.Printf("Reached stream EOF")
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Println(msg.GetBlog())
	}
}
