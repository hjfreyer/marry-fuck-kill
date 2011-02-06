package com.mfk.web.maker.client;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.event.shared.SimpleEventBus;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.RootPanel;
import com.mfk.web.maker.client.event.ImageResultsAvailableEvent;
import com.mfk.web.maker.client.event.QueryUpdatedEvent;
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
    
    Button createButton = new Button("Create");
    RootPanel.get("createButton").add(createButton);

    OutputForm outputForm = new OutputFormImpl(RootPanel.get("outputForm"));

    SimpleEventBus eventBus = new SimpleEventBus();
    
    MfkPresenter presenter = new MfkPresenter(eventBus, ev1, ev2, ev3, 
        createButton, picker, outputForm);
    ImageSearchManager imageSearchManager = new ImageSearchManager(eventBus);
    
    eventBus.addHandler(ImageResultsAvailableEvent.TYPE, presenter);
    eventBus.addHandler(QueryUpdatedEvent.TYPE, imageSearchManager);
  }
}
