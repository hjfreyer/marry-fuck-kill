package com.mfk.web.maker.client.presenter;

import java.util.ArrayList;
import java.util.List;

import com.google.gwt.core.client.JsArray;
import com.google.gwt.search.client.ImageResult;
import com.google.gwt.search.client.ImageSearch;
import com.google.gwt.search.client.Result;
import com.google.gwt.search.client.ResultSetSize;
import com.google.gwt.search.client.SafeSearchValue;
import com.google.gwt.search.client.SearchResultsHandler;
import com.google.gwt.user.client.Timer;
import com.mfk.web.maker.client.event.ImageResultsHandler;
import com.mfk.web.maker.client.model.ImageInfo;

// TODO(mjkelly): See this example:
// https://code.google.com/apis/ajax/playground/#raw_search

public class ImageSearchManager {

  private ImageResultsHandler handler;
  private String currentQuery = "";
  
  private String lastExecutedQuery = "";
  private long lastSearchTime = 0;
  private long lastChangeTime = 0;
  
  static final private long MIN_CHANGE_TIME = 250;
  static final private long MIN_SEARCH_TIME = 1000;
  private Timer searchTimer;

  public void setImageResultsHandler(ImageResultsHandler handler) {
    this.handler = handler;
  }

  public void searchForQuery(final String query) {
    System.out.println("Query changed: " + query);
    currentQuery = query;
    lastChangeTime = System.currentTimeMillis();
    
    if (searchTimer == null) {
      searchTimer = new Timer() {
        public void run() {
          maybeSearch();
        }
      };
      searchTimer.scheduleRepeating(250);
    }
  }
  
  private void maybeSearch() {
    long now = System.currentTimeMillis();
    
    if (!lastExecutedQuery.equals(currentQuery)) {
      if (now - lastChangeTime > MIN_CHANGE_TIME) {
        if (now - lastSearchTime > MIN_SEARCH_TIME) {
          executeSearch(currentQuery);
        }
      }
    }
  }
  
  private void executeSearch(final String query) {
    System.out.println("Searching for: " + query);
    
    lastExecutedQuery = query;
    lastSearchTime = System.currentTimeMillis();
    
    ImageSearch imageSearch = new ImageSearch();
    imageSearch.setSafeSearch(SafeSearchValue.MODERATE);
    imageSearch.setResultSetSize(ResultSetSize.LARGE);
    
    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
      public void onSearchResults(SearchResultsEvent resultsEvent) {
        handleSearchResults(query, resultsEvent.getResults());
      }
    });

    imageSearch.execute(query);
  }

  public void clearState() {
    currentQuery = "";
    lastExecutedQuery = "";
    searchTimer.cancel();
    searchTimer = null;
  }

  private void handleSearchResults(String query,
      JsArray<? extends Result> results) {
    if (!query.equals(lastExecutedQuery)) {
      System.out.println("Discarded results for: " + query);
      return;
    }
    System.out.println("Accepted results for: " + query);

    List<ImageInfo> images = new ArrayList<ImageInfo>();

    for (int i = 0; i < results.length(); i++) {
      ImageResult r = (ImageResult) results.get(i);
      images.add(new ImageInfo(r.getThumbnailUrl(), r.getUnescapedUrl()));
    }

    handler.handleImageResults(query, images);
  }
}
