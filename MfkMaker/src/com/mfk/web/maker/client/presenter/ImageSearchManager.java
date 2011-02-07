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
import com.mfk.web.maker.client.event.ImageResultsHandler;

// TODO(mjkelly): See this example:
// https://code.google.com/apis/ajax/playground/#raw_search

public class ImageSearchManager {
  
  private ImageResultsHandler handler;
  private String currentQuery = "";
  
  public void setImageResultsHandler(ImageResultsHandler handler) {
    this.handler = handler;
  }

  public void searchForQuery(final String query) {
    System.out.println("Searching for: " + query);
            
    ImageSearch imageSearch = new ImageSearch();
    imageSearch.setSafeSearch(SafeSearchValue.MODERATE);
    imageSearch.setResultSetSize(ResultSetSize.LARGE);
    
    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
      public void onSearchResults(SearchResultsEvent resultsEvent) {
        handleSearchResults(query, resultsEvent.getResults());
      }
    });
    
    imageSearch.execute(query);
    currentQuery = query;
  }
  
  public void clearState() {
    currentQuery = "";
  }
  
  private void handleSearchResults(String query, JsArray<? extends Result> results) {
    System.out.println("Got results for " + query);
    
    if (!query.equals(currentQuery)) {
      return;
    }
    
    List<String> urls = new ArrayList<String>();
    
    for (int i = 0; i < results.length(); i++) {
      ImageResult r = (ImageResult) results.get(i);
      urls.add(r.getThumbnailUrl());
    }

    handler.handleImageResults(query, urls);
  }
}
