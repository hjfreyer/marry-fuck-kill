package com.mfk.web.maker.client.presenter;

import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.mfk.web.maker.client.event.EntityPickedHandler;
import com.mfk.web.maker.client.model.EntityInfo;
import com.mfk.web.maker.client.view.EntityView;
import com.mfk.web.maker.client.view.OutputForm;


public class MfkPresenter {

  private final EntityView ev1;
  private final EntityView ev2;
  private final EntityView ev3;
  
  private final OutputForm outputForm;
  
  private EntityPickerPresenter picker;
  
  private EntityInfo entity1 = null;
  private EntityInfo entity2 = null;
  private EntityInfo entity3 = null;
  
  private int currentlyEditing = 0;

  public MfkPresenter(EntityView ev1, EntityView ev2, EntityView ev3,
      OutputForm outputForm, EntityPickerPresenter picker) {
    this.ev1 = ev1;
    this.ev2 = ev2;
    this.ev3 = ev3;
    this.outputForm = outputForm;
    this.picker = picker;
  
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
    
    picker.setEntityPickedHandler(new EntityPickedHandler() {
      @Override
      public void handlePickingCancelled() {
        System.out.println("Picker cancelled");
        currentlyEditing = 0;
      }
      
      @Override
      public void handleEntityPicked(EntityInfo newEntity) {
        System.out.println("Entity picked");

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

        // Enable submit button if all 3 have been changed.
        if (isClickable()) {
          outputForm.setClickable(true);
        }
      }
    });
    
    outputForm.getSubmitButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        submit();
      }
    });
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

    picker.showPicker(entity);
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
