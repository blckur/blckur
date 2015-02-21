library brand_logo_comp;

import 'package:angular/angular.dart' show Component, NgAttr;

@Component(
  selector: 'brand-logo',
  template: '<img src="{{iconUrl}}">',
  cssUrl: 'packages/blckur/components/brand_logo/brand_logo.css'
)
class BrandLogoComp {
  @NgAttr('icon')
  String icon;

  String get iconUrl {
    return '/img/${this.icon}.png';
  }
}
