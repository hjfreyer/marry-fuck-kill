package com.mfk.web.maker.client.event;

import java.util.List;

import com.mfk.web.maker.client.model.ImageInfo;

public interface ImageResultsHandler {
  public void handleImageResults(String query, List<ImageInfo> results);
}
