package com.mfk.web.maker.client.view;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
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
  
  private final Set<String> shownUrls = new HashSet<String>();

  private String selectedUrl = "";
  
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
  public String getSelectedImageUrl() {
    return selectedUrl;
  }

  @Override
  public void setName(String name) {
    searchField.setText(name);
  }

  @Override
  public void addImageUrls(List<String> imageUrls) {
    for (final String imageUrl : imageUrls) {
      if (shownUrls.add(imageUrl)) {
        Image i = new Image(imageUrl);
        imgPanel.add(i);

        i.addClickHandler(new ClickHandler() {
          @Override
          public void onClick(ClickEvent event) {
            selectedUrl = imageUrl;
          }
        });
        
        if (selectedUrl.isEmpty()) {
          selectedUrl = imageUrl;
        }
      }
    }
  }

  public void clearImageUrls() {
    selectedUrl = "";
    shownUrls.clear();
    imgPanel.clear();
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
}
