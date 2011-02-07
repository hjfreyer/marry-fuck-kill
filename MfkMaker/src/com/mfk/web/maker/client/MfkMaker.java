package com.mfk.web.maker.client;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.user.client.ui.RootPanel;
import com.mfk.web.maker.client.presenter.EntityPickerPresenter;
import com.mfk.web.maker.client.presenter.ImageSearchManager;
import com.mfk.web.maker.client.presenter.MfkPresenter;
import com.mfk.web.maker.client.view.EntityPickerView;
import com.mfk.web.maker.client.view.EntityPickerViewImpl;
import com.mfk.web.maker.client.view.EntityView;
import com.mfk.web.maker.client.view.EntityViewImpl;
import com.mfk.web.maker.client.view.OutputForm;
import com.mfk.web.maker.client.view.OutputFormImpl;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkMaker implements EntryPoint {
  /**
   * This is the entry point method.
   */

  public void onModuleLoad() {
    EntityView ev1 = new EntityViewImpl(RootPanel.get("name_1"),
        RootPanel.get("image_1"),
        RootPanel.get("edit_1"));
    EntityView ev2 = new EntityViewImpl(RootPanel.get("name_2"),
        RootPanel.get("image_2"),
        RootPanel.get("edit_2"));
    EntityView ev3 = new EntityViewImpl(RootPanel.get("name_3"),
        RootPanel.get("image_3"),
        RootPanel.get("edit_3"));
    
    EntityPickerView picker = new EntityPickerViewImpl();
    
    OutputForm outputForm = new OutputFormImpl(RootPanel.get("createButton"),
        RootPanel.get("outputForm"));

    ImageSearchManager searchManager = new ImageSearchManager();
    
    EntityPickerPresenter pickerPresenter = 
      new EntityPickerPresenter(picker, searchManager);
    
    new MfkPresenter(ev1, ev2, ev3, outputForm, pickerPresenter);
  }
}
