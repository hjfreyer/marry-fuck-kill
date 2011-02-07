package com.mfk.web.maker.client.event;

import java.util.List;

public interface ImageResultsHandler {
  public void handleImageResults(String query, List<String> resultUrls);
}
