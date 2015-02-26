library brand_logo_comp;

import 'package:blckur/all_aware.dart' as all_aware;

import 'package:angular/angular.dart' show Component, NgAttr;
import 'dart:html' as dom;

@Component(
  selector: 'x-brand-logo',
  template: '<img>'
)
class BrandLogoComp extends all_aware.AllAware {
  @NgAttr('icon')
  String icon;

  @NgAttr('size')
  String size;

  String get iconUrl {
    return '/s/img/${this.icon}.png';
  }

  void onAll(dom.ShadowRoot root) {
    var imgElem = root.querySelector('img');
    if (imgElem != null && this.icon != null) {
      imgElem
        ..style.width = this.size == null ? '80px' : this.size
        ..attributes['src'] = this.iconUrl;
    }
  }
}
