package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/run-x/cloudgrep/graph/generated"
	"github.com/run-x/cloudgrep/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := model.User{
		Name: input.Name,
	}
	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) CreateResource(ctx context.Context, input model.NewResource) (*model.Resource, error) {
	// resource := model.Resource(input)
	resource := model.Resource{
		ID:     input.ID,
		Type:   input.Type,
		Region: input.Region,
		Tags:   mapNewTags(input.Tags),
	}
	err := r.DB.Create(&resource).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("Creating Resource %+v", resource)
	return &resource, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	order := model.Order{
		CustomerName: input.CustomerName,
		OrderAmount:  input.OrderAmount,
		Items:        mapItemsFromInput(input.Items),
	}
	err := r.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]model.User, error) {
	var users []model.User
	r.DB.Find(&users)
	return users, nil
}

func (r *queryResolver) Resources(ctx context.Context) ([]model.Resource, error) {
	var resources []model.Resource
	err := r.DB.Preload("Tags").Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *queryResolver) Orders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := r.DB.Preload("Items").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *queryResolver) Items(ctx context.Context) ([]model.Item, error) {
	var items []model.Item
	r.DB.Find(&items)
	return items, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func mapNewTags(newTags []model.NewTag) []model.Tag {
	var tags []model.Tag
	for _, newTag := range newTags {
		tags = append(tags, model.Tag{
			// ResourceID: 0,
			Key:   newTag.Key,
			Value: newTag.Value,
		})
	}
	return tags
}
func mapItemsFromInput(itemsInput []model.ItemInput) []model.Item {
	var items []model.Item
	for _, itemInput := range itemsInput {
		items = append(items, model.Item{
			ProductCode: itemInput.ProductCode,
			ProductName: itemInput.ProductName,
			Quantity:    itemInput.Quantity,
		})
	}
	return items
}
