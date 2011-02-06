package com.mfk.web.maker.client.view;

import java.util.List;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.event.dom.client.HasKeyUpHandlers;

public interface EntityPickerView {
  public HasClickHandlers getSaveButton();

  public HasClickHandlers getCancelButton();
  
  public HasKeyUpHandlers getSearchTextField();
  
  public String getName();

  public String getSelectedImageUrl();
    
  public void setVisible(boolean visible);
  
  public void setThrob(boolean enabled);
  
  public void setName(String name);
  
  public void addImageUrls(List<String> imageUrl);

  public void clearImageUrls();
}
