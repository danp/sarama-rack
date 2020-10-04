module github.com/danp/sarama-rack

go 1.15

require (
	github.com/Shopify/sarama v1.27.0
	github.com/google/uuid v1.1.2
)

replace github.com/Shopify/sarama => github.com/danp/sarama v1.27.1-0.20201004212245-f5275c211ce9
