package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.event.dom.client.HasKeyUpHandlers;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.Panel;
import com.google.gwt.user.client.ui.TextBox;

public class EntityPickerViewImpl implements EntityPickerView {

  private final DialogBox dialog = new DialogBox();
  private Image autoThrobber = new Image("/s/loading.gif");

  private final TextBox searchField = new TextBox();
  private final Button saveButton = new Button("Save");
  private final Button cancelButton = new Button("Cancel");

  private final Image imagePreview = new Image();
  private final Panel imageChoices;

  public EntityPickerViewImpl() {
    Panel dialogPanel = new FlowPanel();
    dialogPanel.setStyleName("pickerBody");

    Panel rightside = makeSubPanel(dialogPanel, "rightSide");
    imageChoices = makeSubPanel(dialogPanel, "imageChoices");
    Panel submitBar = makeSubPanel(dialogPanel, "submitBar");

    rightside.add(new Label("Item name:"));
    rightside.add(searchField);
    rightside.add(autoThrobber);

    Panel primaryImgBox = makeSubPanel(rightside, "primaryImgBox");
    primaryImgBox.add(imagePreview);

    submitBar.add(cancelButton);
    submitBar.add(saveButton);
    
    cancelButton.setStyleName("button");
    cancelButton.addStyleName("clickable");
    saveButton.setStyleName("button");

    imagePreview.setStyleName("imagePreview");

    dialog.setText("Pick your poison");
    dialog.addStyleName("picker");
    dialog.setGlassEnabled(true);
    dialog.setWidget(dialogPanel);
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

    imageChoices.add(imageView.getImage());

    return imageView;
  }

  @Override
  public void setThrob(boolean enabled) {
    if (enabled) {
      autoThrobber.removeStyleName("hidden");
    } else {
      autoThrobber.addStyleName("hidden");
    }
  }

  @Override
  public void setVisible(boolean visible) {
    if (visible) {
      dialog.center();
      dialog.show();
      searchField.setFocus(true);
    } else {
      dialog.hide();
    }
  }

  @Override
  public void clearImages() {
    imageChoices.clear();
  }

  @Override
  public void showPreview(String imageUrl) {
    imagePreview.setUrl(imageUrl);
    imagePreview.setVisible(true);
  }

  @Override
  public void setSaveable(boolean saveable) {
    if (saveable) {
      saveButton.addStyleName("clickable");
    } else {
      saveButton.removeStyleName("clickable");
    }    
  }
  
  @Override
  public void clearPreview() {
    imagePreview.setVisible(false);
  }

  private Panel makeSubPanel(Panel parent, String styleName) {
    Panel result = new FlowPanel();
    result.setStyleName(styleName);
    parent.add(result);

    return result;
  }
}
