package com.mfk.web.maker.client.event;

import com.mfk.web.maker.client.model.EntityInfo;

public interface EntityPickedHandler {

  public void handleEntityPicked(EntityInfo entity);
   
  public void handlePickingCancelled();
}
