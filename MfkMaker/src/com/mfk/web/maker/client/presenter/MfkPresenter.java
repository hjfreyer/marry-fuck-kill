package com.mfk.web.maker.client.presenter;

import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.event.shared.HasHandlers;
import com.google.gwt.user.client.Window;
import com.mfk.web.maker.client.event.ImageResultsAvailableEvent;
import com.mfk.web.maker.client.event.QueryUpdatedEvent;
import com.mfk.web.maker.client.model.EntityInfo;
import com.mfk.web.maker.client.view.EntityPickerView;
import com.mfk.web.maker.client.view.EntityView;
import com.mfk.web.maker.client.view.OutputForm;


public class MfkPresenter implements ImageResultsAvailableEvent.Handler {

  private final HasHandlers eventBus;

  private final EntityInfo default1 = 
    new EntityInfo("Treehouse 1", "/s/treehouse-1.jpeg");
  private final EntityInfo default2 = 
    new EntityInfo("Treehouse 2", "/s/treehouse-2.jpeg");
  private final EntityInfo default3 = 
    new EntityInfo("Treehouse 3", "/s/treehouse-3.jpeg");
  
  private final EntityView ev1;
  private final EntityView ev2;
  private final EntityView ev3;
  
  private final EntityPickerView picker;
  
  private final HasClickHandlers createButton;

  private final OutputForm outputForm;
  
  private EntityView currentlyEditing = null;
  private String lastEnteredQuery = "";
  
  public MfkPresenter(HasHandlers eventBus, EntityView ev1, EntityView ev2,
      EntityView ev3, HasClickHandlers createButton, EntityPickerView picker,
      OutputForm outputForm) {
    this.eventBus = eventBus;
    this.ev1 = ev1;
    this.ev2 = ev2;
    this.ev3 = ev3;
    this.createButton = createButton;
    this.picker = picker;
    this.outputForm = outputForm;

    init();
  }

  private void init() {
    ev1.getEditButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        editEntity(ev1);
      }
    });

    ev2.getEditButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        editEntity(ev2);
      }
    });
    
    ev3.getEditButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        editEntity(ev3);
      }
    });
    
    picker.getSaveButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        pickerSave();
      }
    });

    picker.getCancelButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        pickerCancel();
      }
    });
    
    createButton.addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        submit();
      }
    });
        
    picker.getSearchTextField().addKeyUpHandler(new KeyUpHandler() {
      @Override
      public void onKeyUp(KeyUpEvent event) {
        onQueryUpdated();
      }
    });
    
    picker.setVisible(false);
    ev1.setEntityInfo(default1);
    ev2.setEntityInfo(default2);
    ev3.setEntityInfo(default3);
  }
  
  public void onQueryUpdated() {
    if (!picker.getName().equals(lastEnteredQuery)) {
      picker.setThrob(true);
      eventBus.fireEvent(new QueryUpdatedEvent(picker.getName()));
      lastEnteredQuery = picker.getName();
    }
  }

  @Override
  public void handleImageResults(ImageResultsAvailableEvent event) {
    System.out.println("Image results available for " + event.query);
    
    if (event.newResultSet) {
      picker.clearImageUrls();
    }

    picker.setThrob(false);
    picker.addImageUrls(event.resultUrls);
  }

  public void editEntity(EntityView entityView) {
    currentlyEditing = entityView;

    picker.setName(entityView.getEntityInfo().getName());    
    picker.clearImageUrls();
    picker.setVisible(true);
    
    eventBus.fireEvent(new QueryUpdatedEvent(picker.getName()));
  }
  
  private void pickerSave() {
    EntityInfo newEntity = new EntityInfo(picker.getName(),
                                          picker.getSelectedImageUrl());
    currentlyEditing.setEntityInfo(newEntity);

    currentlyEditing = null;
    picker.setVisible(false); 
    
    // This is a hacky way of telling the ImageSearchManager that 
    // we are done with this searching session and to clear its state.
    // Make this a separate event?
    eventBus.fireEvent(new QueryUpdatedEvent(""));
  }

  private void pickerCancel() {
    currentlyEditing = null;
    picker.setVisible(false);
    // This is a hacky way of telling the ImageSearchManager that 
    // we are done with this searching session and to clear its state.
    // Make this a separate event?
    eventBus.fireEvent(new QueryUpdatedEvent(""));
  }
  
  private void submit() {
    if (!ev1.getEntityInfo().getImageUrl().startsWith("http://images.google.com") ||
        !ev2.getEntityInfo().getImageUrl().startsWith("http://images.google.com") ||
        !ev3.getEntityInfo().getImageUrl().startsWith("http://images.google.com")) {    
      Window.alert("You must change all the items from their defaults.");
      return;
    }

    outputForm.setEntity1(ev1.getEntityInfo());
    outputForm.setEntity2(ev2.getEntityInfo());
    outputForm.setEntity3(ev3.getEntityInfo());
    
    outputForm.submit();
  }
}
