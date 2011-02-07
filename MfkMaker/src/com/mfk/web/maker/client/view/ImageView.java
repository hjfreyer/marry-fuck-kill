package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;

public interface ImageView {

  public HasClickHandlers getClickable();

  public void setImageUrl(String url);

  public void setSelected(boolean selected);

}
