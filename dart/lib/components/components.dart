library components;

import 'package:blckur/components/account_item/account_item.dart'
  as account_item;
import 'package:blckur/components/brand_logo/brand_logo.dart'
  as brand_logo;
import 'package:blckur/components/feed/feed.dart'
  as feed;
import 'package:blckur/components/test/test.dart'
  as test;

import 'package:angular/angular.dart' as ng;

class ComponentsMod extends ng.Module {
  ComponentsMod() {
    this.bind(account_item.AccountItemComp);
    this.bind(brand_logo.BrandLogoComp);
    this.bind(feed.FeedComp);
    this.bind(test.TestComp);
  }
}
