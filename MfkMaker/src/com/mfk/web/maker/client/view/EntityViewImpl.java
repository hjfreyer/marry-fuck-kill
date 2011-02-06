package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.Panel;
import com.mfk.web.maker.client.model.EntityInfo;

public class EntityViewImpl implements EntityView {

  private Label name = new Label();
  private Image img = new Image();
  private Button editButton = new Button("Edit");
  
  public EntityViewImpl(Panel name, Panel img, Panel edit) {
    name.add(this.name);
    img.add(this.img);
    edit.add(editButton);

    editButton.removeStyleName("gwt-Button");
    editButton.addStyleName("button");
    editButton.addStyleName("clickable");
  }

  @Override
  public HasClickHandlers getEditButton() {
    return editButton;
  }
  
  @Override
  public EntityInfo getEntityInfo() {
    return new EntityInfo(name.getText(), img.getUrl());
  }

  @Override
  public void setEntityInfo(EntityInfo info) {
    name.setText(info.getName());
    img.setUrl(info.getImageUrl());
  }
}
