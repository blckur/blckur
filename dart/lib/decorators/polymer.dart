library polymer_dec;

import 'package:angular/angular.dart' show Decorator;
import 'package:angular/angular.dart' as ng;
import 'package:observe/observe.dart' as obsrv;
import 'package:template_binding/template_binding.dart' as tmpl_bnd;
import 'dart:html' as dom;

typedef dynamic _Nullary();

@Decorator(
  selector: '[py-*]',
  exportExpressionAttrs: const [
    'py-activateEvent',
    'py-active',
    'py-activeClass',
    'py-allowOverflow',
    'py-alt',
    'py-animated',
    'py-auto',
    'py-autoCloseDisabled',
    'py-autoFocusDisabled',
    'py-autoSaveDisabled',
    'py-backdrop',
    'py-baseClass',
    'py-body',
    'py-bottomJustify',
    'py-checked',
    'py-closeAttribute',
    'py-closeSelector',
    'py-closedClass',
    'py-closedIcon',
    'py-committedValue',
    'py-completeEventName',
    'py-completed',
    'py-condensedHeaderHeight',
    'py-condenses',
    'py-contentType',
    'py-data',
    'py-defaultSelected',
    'py-direction',
    'py-disableDrag',
    'py-disableSwipe',
    'py-disabled',
    'py-drawerWidth',
    'py-duration',
    'py-editable',
    'py-error',
    'py-excludedLocalNames',
    'py-fade',
    'py-fill',
    'py-fixed',
    'py-fixedSize',
    'py-floatingLabel',
    'py-focused',
    'py-for',
    'py-forceNarrow',
    'py-grid',
    'py-groups',
    'py-halign',
    'py-handleAs',
    'py-headerHeight',
    'py-headers',
    'py-heading',
    'py-height',
    'py-hideScrollButton',
    'py-horizontal',
    'py-icon',
    'py-iconSize',
    'py-icons',
    'py-id',
    'py-immediateValue',
    'py-indeterminate',
    'py-initialOpacity',
    'py-isInvalid',
    'py-itemsSelector',
    'py-justify',
    'py-keepCondensedHeader',
    'py-keys',
    'py-label',
    'py-labelVisible',
    'py-lastSelected',
    'py-layered',
    'py-list',
    'py-loading',
    'py-locked',
    'py-lowerThreshold',
    'py-lowerTriggered',
    'py-max',
    'py-maxRows',
    'py-method',
    'py-middleJustify',
    'py-min',
    'py-minSize',
    'py-mini',
    'py-mode',
    'py-multi',
    'py-name',
    'py-narrow',
    'py-noDissolve',
    'py-noReveal',
    'py-noarrow',
    'py-nobar',
    'py-noink',
    'py-noslide',
    'py-notap',
    'py-offsetX',
    'py-offsetY',
    'py-opacityDecayVelocity',
    'py-opened',
    'py-openedClass',
    'py-openedIcon',
    'py-orient',
    'py-params',
    'py-pin',
    'py-placeholder',
    'py-position',
    'py-preload',
    'py-pressed',
    'py-preventInvalidInput',
    'py-progress',
    'py-query',
    'py-queryMatches',
    'py-raised',
    'py-ratio',
    'py-recenteringTouch',
    'py-ref',
    'py-relatedTarget',
    'py-response',
    'py-responsiveWidth',
    'py-rightDrawer',
    'py-rows',
    'py-runwayFactor',
    'py-scopeClass',
    'py-scrollAwayTopbar',
    'py-scrollTarget',
    'py-scrollable',
    'py-secondaryProgress',
    'py-selected',
    'py-selectedAttribute',
    'py-selectedClass',
    'py-selectedIndex',
    'py-selectedItem',
    'py-selectedModel',
    'py-selectedProperty',
    'py-selection',
    'py-selectionEnabled',
    'py-shadow',
    'py-show',
    'py-sizing',
    'py-sizingTarget',
    'py-snaps',
    'py-src',
    'py-step',
    'py-swipeDisabled',
    'py-tallClass',
    'py-target',
    'py-text',
    'py-tipAttribute',
    'py-toggle',
    'py-toggles',
    'py-transition',
    'py-transitionProperty',
    'py-transitionType',
    'py-transitions',
    'py-type',
    'py-upperThreshold',
    'py-upperTriggered',
    'py-url',
    'py-useRaw',
    'py-valign',
    'py-value',
    'py-valueattr',
    'py-width',
    'py-withCredentials',
    'py-z',
  ]
)
class PolymerDec {
  PolymerDec(dom.Node node, ng.Parser parser, ng.Scope scope) {
    dom.Element element = node;

    element.attributes.forEach((key, value) {
      if (!key.startsWith('py-')) {
        return;
      }
      var bindVal = new ValueBindable();
      var expression = parser(value);

      tmpl_bnd.nodeBind(node).bind(key.substring(3), bindVal);

      if (expression.isAssignable) {
        bindVal.onChange = (value) {
          expression.assign(scope.context, value);
        };
      }

      scope.watch(value, bindVal.update);
    });
  }
}

class ValueBindable implements obsrv.Bindable {
  Function callback;
  Function onChange;

  var _value;
  void set value(dynamic value) {
    this._value = value;
    if (this.onChange != null) {
      this.onChange(this._value);
    }
  }
  dynamic get value {
    return this._value;
  }

  void close() {
    this.callback = null;
  }

  dynamic open(Function callback) {
    this.callback = callback;
    return this._value;
  }

  void update(dynamic newValue, dynamic oldValue) {
    this._value = newValue;
    if (this.callback != null) {
      if (this.callback is _Nullary) {
        this.callback();
      }
      else {
        this.callback(this._value);
      }
    }
  }

  void deliver() {}
}
