package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.event.dom.client.LoadEvent;
import com.google.gwt.event.dom.client.LoadHandler;
import com.google.gwt.user.client.ui.Image;

public class ImageViewImpl implements ImageView {
  
  private final Image image = new Image();
  
  @Override
  public HasClickHandlers getClickable() {
    return image;
  }
  
  @Override
  public void setImageUrl(String url) {
    image.setVisible(false);
    image.setUrl(url);
    image.addLoadHandler(new LoadHandler() {
      @Override
      public void onLoad(LoadEvent event) {
        image.setVisible(true);
      }
    });
  }
  
  @Override
  public void setSelected(boolean selected) {

  }
  
  public Image getImage() {
    return image;
  }
}
