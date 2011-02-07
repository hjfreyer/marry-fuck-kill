package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.event.dom.client.HasKeyUpHandlers;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Panel;
import com.google.gwt.user.client.ui.TextBox;

public class EntityPickerViewImpl implements EntityPickerView {

  private final DialogBox dialog = new DialogBox();
  private final Panel panel = new FlowPanel();

  private Image autoThrobber = new Image("/s/loading.gif");
  private final Panel imgPanel = new FlowPanel();

  private final TextBox searchField = new TextBox();
  private final Button saveButton = new Button("Save");
  private final Button cancelButton = new Button("Cancel");
  
  public EntityPickerViewImpl() {
    panel.add(searchField);
    panel.add(cancelButton);
    panel.add(saveButton);
    panel.add(autoThrobber);
    panel.add(imgPanel);
    
    panel.setWidth("600px");
    panel.setHeight("500px");

    dialog.setText("Pick your poison");
    dialog.setGlassEnabled(true);
    dialog.setWidget(panel);
    dialog.center();
  }

  @Override
  public HasClickHandlers getCancelButton() {
    return cancelButton;
  }

  @Override
  public HasClickHandlers getSaveButton() {
    return saveButton;
  }

  @Override
  public HasKeyUpHandlers getSearchTextField() {
    return searchField;
  }

  @Override
  public String getName() {
    return searchField.getText();
  }

  @Override
  public void setName(String name) {
    searchField.setText(name);
  }

  @Override
  public ImageView addImageView() {
    ImageViewImpl imageView = new ImageViewImpl();
    
    imgPanel.add(imageView.getImage());
    
    return imageView;
  }
    
  @Override
  public void setThrob(boolean enabled) {
    autoThrobber.setVisible(enabled);    
  }

  @Override
  public void setVisible(boolean visible) {
    if (visible) {
      dialog.show();
      searchField.setFocus(true);
    } else {
      dialog.hide();
    }
  }

  @Override
  public void clearImages() {
    imgPanel.clear();
  }

  @Override
  public void showPreview(String imageUrl) {
    
  }
  
  @Override
  public void clearPreview() {
    // TODO Auto-generated method stub
    
  }
}
