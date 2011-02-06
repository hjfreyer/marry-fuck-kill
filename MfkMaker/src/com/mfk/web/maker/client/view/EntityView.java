package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.mfk.web.maker.client.model.EntityInfo;

public interface EntityView {
  public HasClickHandlers getEditButton();
  
  public EntityInfo getEntityInfo();

  public void setEntityInfo(EntityInfo info);
}
