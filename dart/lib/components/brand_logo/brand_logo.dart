library brand_logo_comp;

import 'package:blckur/constants.dart';
import 'package:blckur/all_aware.dart' as all_aware;

import 'package:angular/angular.dart' show Component, NgAttr;
import 'dart:html' as dom;

@Component(
  selector: 'x-brand-logo',
  template: '<img>'
)
class BrandLogoComp extends all_aware.AllAware {
  dom.Element imgElem;

  var _icon;
  @NgAttr('icon')
  String get icon {
    return this._icon;
  }
  void set icon(String val) {
    if (val != this._icon) {
      this._icon = val;
      this._updateIcon();
      this._updateSize();
    }
  }

  var _size;
  @NgAttr('size')
  String get size {
    return this._size;
  }
  void set size(String val) {
    if (val != this._size) {
      this._size = val;
      this._updateSize();
    }
  }

  String get iconUrl {
    var icon = icons[this.icon];

    if (icon == null) {
      icon = "";
    }

    return icon['url'];
  }

  void _updateIcon() {
    if (this.imgElem != null && this.icon != null) {
      this.imgElem.attributes['src'] = this.iconUrl;
    }
  }

  void _updateSize() {
    var icon = icons[this.icon];

    if (this.imgElem != null) {
      this.imgElem.style
        ..maxWidth = this.size == null ? '80px' : this.size
        ..maxHeight = this.size == null ? '58px' : this.size
        ..margin = '0 auto';

      if (icon != null && icon['style'] == NARROW) {
        this.imgElem.style
          ..width = 'auto'
          ..height = '100%';
      } else {
        this.imgElem.style
          ..width = '100%'
          ..height = 'auto';
      }
    }
  }

  void onAll(dom.ShadowRoot root) {
    this.imgElem = root.querySelector('img');
    this._updateIcon();
    this._updateSize();
  }
}
