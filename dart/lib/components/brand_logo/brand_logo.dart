library brand_logo_comp;

import 'package:angular/angular.dart' show Component, NgAttr;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;

@Component(
  selector: 'brand-logo',
  template: '<img>'
)
class BrandLogoComp implements ng.ShadowRootAware {
  @NgAttr('icon')
  String icon;

  @NgAttr('size')
  String size;

  String get iconUrl {
    return '/img/${this.icon}.png';
  }

  void onShadowRoot(dom.ShadowRoot root) {
    root.querySelector('img')
      ..style.width = this.size == null ? '80px' : this.size
      ..attributes['src'] = this.iconUrl;
  }
}
