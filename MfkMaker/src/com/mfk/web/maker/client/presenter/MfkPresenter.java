package com.mfk.web.maker.client.presenter;

import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.event.shared.HasHandlers;
import com.mfk.web.maker.client.event.ImageResultsAvailableEvent;
import com.mfk.web.maker.client.event.QueryUpdatedEvent;
import com.mfk.web.maker.client.model.EntityInfo;
import com.mfk.web.maker.client.view.EntityPickerView;
import com.mfk.web.maker.client.view.EntityView;
import com.mfk.web.maker.client.view.OutputForm;


public class MfkPresenter implements ImageResultsAvailableEvent.Handler {

  private final HasHandlers eventBus;
  
  private final EntityView ev1;
  private final EntityView ev2;
  private final EntityView ev3;
  
  private final EntityPickerView picker;
  
  private final OutputForm outputForm;
  
  private EntityInfo entity1 = null;
  private EntityInfo entity2 = null;
  private EntityInfo entity3 = null;
  
  private int currentlyEditing = 0;
  private String lastEnteredQuery = "";
  
  public MfkPresenter(HasHandlers eventBus, EntityView ev1, EntityView ev2,
      EntityView ev3, EntityPickerView picker,
      OutputForm outputForm) {
    this.eventBus = eventBus;
    this.ev1 = ev1;
    this.ev2 = ev2;
    this.ev3 = ev3;
    this.picker = picker;
    this.outputForm = outputForm;

    init();
  }

  private void init() {
    ev1.getEditButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        editEntity(1);
      }
    });

    ev2.getEditButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        editEntity(2);
      }
    });
    
    ev3.getEditButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        editEntity(3);
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
    
    outputForm.getSubmitButton().addClickHandler(new ClickHandler() {
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

  public void editEntity(int entityIndex) {
    currentlyEditing = entityIndex;
    
    EntityInfo entity = null;
    
    if (currentlyEditing == 1) {
      entity = entity1;
    } else if (currentlyEditing == 2) {
      entity = entity2;
    } else if (currentlyEditing == 3) {
      entity = entity3;
    }

    picker.clearImageUrls();
    if (entity == null) {
      picker.setName("");
    } else {
      picker.setName(entity.name);    
      eventBus.fireEvent(new QueryUpdatedEvent(entity.name));
    }
    picker.setVisible(true);
  }
  
  private void pickerSave() {
    if (picker.getName().isEmpty() || picker.getSelectedImageUrl().isEmpty()) {
      pickerCancel();
      return;
    }

    EntityInfo newEntity = new EntityInfo(picker.getName(),
                                          picker.getSelectedImageUrl());
 
    if (currentlyEditing == 1) {
      entity1 = newEntity;
      ev1.showEntity(newEntity);
    } else if (currentlyEditing == 2) {
      entity2 = newEntity;
      ev2.showEntity(newEntity);
    } else if (currentlyEditing == 3) {
      entity3 = newEntity;
      ev3.showEntity(newEntity);
    }

    currentlyEditing = 0;
    picker.setVisible(false); 
    
    // Enable submit button if all 3 have been changed.
    if (isClickable()) {
      outputForm.setClickable(true);
    }
    
    // This is a hacky way of telling the ImageSearchManager that 
    // we are done with this searching session and to clear its state.
    // Make this a separate event?
    eventBus.fireEvent(new QueryUpdatedEvent(""));
  }

  private void pickerCancel() {
    currentlyEditing = 0;
    picker.setVisible(false);
    // This is a hacky way of telling the ImageSearchManager that 
    // we are done with this searching session and to clear its state.
    // Make this a separate event?
    eventBus.fireEvent(new QueryUpdatedEvent(""));
  }
  
  private boolean isClickable() {
    return entity1 != null && entity2 != null && entity3 != null;
  }
  
  private void submit() {
    if (!isClickable()) {
      return;
    }
    
    outputForm.setEntity1(entity1);
    outputForm.setEntity2(entity2);
    outputForm.setEntity3(entity3);
    
    outputForm.submit();
  }
}
