library components;

import 'package:blckur/components/account/account.dart'
  as account;
import 'package:blckur/components/account_filters/account_filters.dart'
  as account_filters;
import 'package:blckur/components/accounts/accounts.dart'
  as accounts;
import 'package:blckur/components/alerts/alerts.dart'
  as alerts;
import 'package:blckur/components/auth/auth.dart'
  as auth;
import 'package:blckur/components/brand_logo/brand_logo.dart'
  as brand_logo;
import 'package:blckur/components/feed/feed.dart'
  as feed;
import 'package:blckur/components/get_started/get_started.dart'
  as get_started;
import 'package:blckur/components/input/input.dart'
  as input;
import 'package:blckur/components/user/user.dart'
  as user;
import 'package:blckur/components/notification/notification.dart'
  as notification;
import 'package:blckur/components/notifications/notifications.dart'
  as notifications;
import 'package:blckur/components/recaptcha/recaptcha.dart'
  as recaptcha;

import 'package:angular/angular.dart' as ng;

class ComponentsMod extends ng.Module {
  ComponentsMod() {
    this.bind(account.AccountComp);
    this.bind(account_filters.AccountFiltersComp);
    this.bind(accounts.AccountsComp);
    this.bind(alerts.AlertsComp);
    this.bind(auth.AuthComp);
    this.bind(brand_logo.BrandLogoComp);
    this.bind(feed.FeedComp);
    this.bind(get_started.GetStartedComp);
    this.bind(input.InputComp);
    this.bind(user.UserComp);
    this.bind(notification.NotificationComp);
    this.bind(notifications.NotificationsComp);
    this.bind(recaptcha.RecaptchaComp);
  }
}
