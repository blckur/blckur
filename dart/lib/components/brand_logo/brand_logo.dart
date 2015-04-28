library brand_logo_comp;

import 'package:blckur/all_aware.dart' as all_aware;

import 'package:angular/angular.dart' show Component, NgAttr;
import 'dart:html' as dom;

const WIDE = 0;
const NARROW = 1;

// Used for static hash versioning search and replace
var icons = {
  'digitalocean': {
    'url': '/s/img/digitalocean.png',
    'style': WIDE,
  },
  'github': {
    'url': '/s/img/github.png',
    'style': WIDE,
   },
  'gmail': {
    'url': '/s/img/gmail.png',
    'style': WIDE,
   },
  'hackernews': {
    'url': '/s/img/hackernews.png',
    'style': NARROW,
   },
  'hipchat': {
    'url': '/s/img/hipchat.png',
    'style': WIDE,
   },
  'stripe': {
    'url': '/s/img/stripe.png',
    'style': WIDE,
   },
  'twitter': {
    'url': '/s/img/twitter.png',
    'style': WIDE,
   },
};

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
