package main

import (
  "fmt"
  webpush "github.com/sherclockholmes/webpush-go"
  "github.com/djavorszky/ddn/common/logger"
)

var userSubscriptions []webpush.Subscription

// sends a notification to a certain user's subscribed endpoints (Chrome, Firefox, etc.)
func sendUserNotifications(subscriber string, message string) error {
 
  userSubscriptions, err := db.FetchUserPushSubscriptions(subscriber);
 
  if err != nil {
    logger.Error("Error:", err)
    return err
  }

  for _, subscription := range userSubscriptions {
    _, err := webpush.SendNotification([]byte(message), &subscription, &webpush.Options{
      Subscriber:      "clouddb@liferay.com",
      VAPIDPrivateKey: config.VAPIDPrivateKey,
    })

    if err != nil {
      logger.Error(fmt.Sprintf("Failed sending notification for user %v to endpoint %v", subscriber, subscription.Endpoint))
      logger.Error(fmt.Sprintf("Returned eror: %v", err))
      return err
    }
  }

  return nil
  
}