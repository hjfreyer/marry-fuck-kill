package com.mfk.web.maker.client.presenter;

import java.util.List;

import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.mfk.web.maker.client.event.EntityPickedHandler;
import com.mfk.web.maker.client.event.ImageResultsHandler;
import com.mfk.web.maker.client.model.EntityInfo;
import com.mfk.web.maker.client.model.ImageInfo;
import com.mfk.web.maker.client.view.EntityPickerView;
import com.mfk.web.maker.client.view.ImageView;

public class EntityPickerPresenter implements ImageResultsHandler {

  private final EntityPickerView view;

  private final ImageSearchManager searchManager;

  private EntityPickedHandler pickedHandler = null;

  private String enteredName = "";
  private String previewedUrl = "";
  private String resultsQuery = "";
  private String originalImageUrl = "";

  private ImageView currentlySelectedImage = null;
  private boolean saveable = false;


  public EntityPickerPresenter(EntityPickerView view,
      ImageSearchManager searchManager) {
    this.view = view;
    this.searchManager = searchManager;

    init();
  }

  private void init() {
    view.getCancelButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        cancel();
      }
    });

    view.getSaveButton().addClickHandler(new ClickHandler() {
      @Override
      public void onClick(ClickEvent event) {
        if (saveable) {
          save();
        }
      }
    });

    view.getSearchTextField().addKeyUpHandler(new KeyUpHandler() {
      @Override
      public void onKeyUp(KeyUpEvent event) {
        if (!view.getName().equals(enteredName)) {
          enteredName = view.getName();
          updateQuery();
        }
      }
    });

    searchManager.setImageResultsHandler(this);

    view.setVisible(false);
  }

  public void setEntityPickedHandler(EntityPickedHandler pickedHandler) {
    this.pickedHandler = pickedHandler;
  }

  public void save() {
    view.setVisible(false);

    if (enteredName.isEmpty() || previewedUrl.isEmpty()) {
      pickedHandler.handlePickingCancelled();
    } else {
      pickedHandler.handleEntityPicked(
          new EntityInfo(enteredName, previewedUrl, originalImageUrl));
    }
    
    searchManager.clearState();
  }

  public void cancel() {
    view.setVisible(false);
    pickedHandler.handlePickingCancelled();
    searchManager.clearState();
  }

  public void updateQuery() {
    if (enteredName.isEmpty()) {
      view.clearImages();
      view.clearPreview();
      searchManager.searchForQuery("");
      
      saveable = false;
      view.setSaveable(false);
    } else if (!enteredName.equals(resultsQuery)) {
      view.setThrob(true);
      searchManager.searchForQuery(enteredName);
      
      saveable = true;
      view.setSaveable(true);
    } else {
      view.setThrob(false);
      searchManager.searchForQuery(enteredName);
      
      saveable = true;
      view.setSaveable(true);
    }
  }

  // Null if we have no starting entity.
  public void showPicker(EntityInfo entity) {
    if (entity == null) {
      enteredName = "";
      previewedUrl = "";
      originalImageUrl = "";

      view.setName("");
      view.clearPreview();
      view.clearImages();

      view.setThrob(false);
      saveable = false;
      view.setSaveable(false);
    } else {
      enteredName = entity.name;
      previewedUrl = entity.imageUrl;
      originalImageUrl = entity.originalImageUrl;

      view.setName(enteredName);
      view.showPreview(previewedUrl);
      view.clearImages();

      view.setThrob(true);
      searchManager.searchForQuery(enteredName);
      saveable = true;
      view.setSaveable(true);
    }

    view.setVisible(true);
  }

  @Override
  public void handleImageResults(String query, List<ImageInfo> results) {
    resultsQuery = query;
    view.clearPreview();
    view.clearImages();
    previewedUrl = "";

    for (final ImageInfo result : results) {
      final ImageView img = view.addImageView();
      img.setImageUrl(result.url);
      img.setSelected(false);

      if (previewedUrl.isEmpty()) {
        previewedUrl = result.url;
        originalImageUrl = result.originalUrl;
        view.showPreview(result.url);
      }

      if (previewedUrl.equals(result.url)) {
        img.setSelected(true);
        currentlySelectedImage = img;
      }

      img.getClickable().addClickHandler(new ClickHandler() {
        @Override
        public void onClick(ClickEvent event) {
          previewedUrl = result.url;
          originalImageUrl = result.originalUrl;

          if (currentlySelectedImage != null) {
            currentlySelectedImage.setSelected(false);
          }

          currentlySelectedImage = img;

          img.setSelected(true);

          view.showPreview(result.url);
        }
      });
    }

    view.setThrob(false);
  }
}
