package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password        *string            `json:"password"   validate:"required,min=6"`
	Email           *string            `json:"email"      validate:"email,required"`
	Phone           *string            `json:"phone"      validate:"required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `josn:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updtaed_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Image        *string            `json:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price"  bson:"price"`
	Rating       *uint              `json:"rating" bson:"rating"`
	Image        *string            `json:"image"  bson:"image"`
}

type Address struct {
	Address_id primitive.ObjectID `bson:"_id"`
	House      *string            `json:"house_name" bson:"house_name"`
	Street     *string            `json:"street_name" bson:"street_name"`
	City       *string            `json:"city_name" bson:"city_name"`
	Pincode    *string            `json:"pin_code" bson:"pin_code"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list"  bson:"order_list"`
	Orderered_At   time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount"    bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod"     bson:"cod"`
}
