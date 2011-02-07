package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.mfk.web.maker.client.model.EntityInfo;

public interface OutputForm {
  public HasClickHandlers getSubmitButton();

  public void setClickable(boolean clickable);

  public void setEntity1(EntityInfo info);

  public void setEntity2(EntityInfo info);

  public void setEntity3(EntityInfo info);

  public void submit();
}
