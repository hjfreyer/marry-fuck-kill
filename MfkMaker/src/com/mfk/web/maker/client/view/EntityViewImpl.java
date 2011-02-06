package com.mfk.web.maker.client.view;

import com.google.gwt.dom.client.Style.Unit;
import com.google.gwt.dom.client.Style.Visibility;
import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.Panel;
import com.mfk.web.maker.client.model.EntityInfo;

public class EntityViewImpl implements EntityView {

  private final Label name = new Label();
  private final Image img = new Image();
  private final Button editButton = new Button("Edit");
  
  public EntityViewImpl(Panel namePanel, Panel imgPanel, Panel editPanel) {
    namePanel.add(name);
    imgPanel.add(img);
    editPanel.add(editButton);

    editButton.removeStyleName("gwt-Button");
    editButton.addStyleName("button");
    editButton.addStyleName("clickable");

    name.setText("?");
    img.setUrl("/s/mfk.png");

    name.getElement().getStyle().setVisibility(Visibility.HIDDEN);
    img.getElement().getStyle().setOpacity(0.3);
    img.getElement().getStyle().setMarginTop(-20, Unit.PX);
  }

  @Override
  public HasClickHandlers getEditButton() {
    return editButton;
  }
  
  @Override
  public void showEntity(EntityInfo info) {
    name.getElement().getStyle().clearVisibility();
    img.getElement().getStyle().clearOpacity();
    img.getElement().getStyle().clearMarginTop();
    
    name.setText(info.getName());
    img.setUrl(info.getImageUrl());
  }
}