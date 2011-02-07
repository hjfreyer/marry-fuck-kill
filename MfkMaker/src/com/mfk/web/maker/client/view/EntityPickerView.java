package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.event.dom.client.HasKeyUpHandlers;

public interface EntityPickerView {
  public HasClickHandlers getSaveButton();

  public HasClickHandlers getCancelButton();
  
  public HasKeyUpHandlers getSearchTextField();
  
  public String getName();
    
  public void setVisible(boolean visible);
  
  public void setThrob(boolean enabled);
  
  public void setName(String name);
  
  public ImageView addImageView();

  public void clearImages();

  public void showPreview(String imageUrl);

  public void clearPreview();
}
