package com.mfk.web.maker.client.presenter;

import java.util.ArrayList;
import java.util.List;

import com.google.gwt.core.client.JsArray;
import com.google.gwt.event.shared.HasHandlers;
import com.google.gwt.search.client.ImageResult;
import com.google.gwt.search.client.ImageSearch;
import com.google.gwt.search.client.Result;
import com.google.gwt.search.client.ResultSetSize;
import com.google.gwt.search.client.SafeSearchValue;
import com.google.gwt.search.client.SearchResultsHandler;
import com.mfk.web.maker.client.event.ImageResultsAvailableEvent;
import com.mfk.web.maker.client.event.QueryUpdatedEvent;

// TODO(mjkelly): See this example:
// https://code.google.com/apis/ajax/playground/#raw_search

public class ImageSearchManager implements QueryUpdatedEvent.Handler {
  
  private final HasHandlers eventBus;
    
  private String currentQuery = "";
  private String lastReceivedQuery = "";
  
  public ImageSearchManager(HasHandlers eventBus) {
    this.eventBus = eventBus;
  }

  @Override
  public void handleNewQuery(final QueryUpdatedEvent event) {
    System.out.println("Query changed to: " + event.query);
            
    ImageSearch imageSearch = new ImageSearch();
    imageSearch.setSafeSearch(SafeSearchValue.MODERATE);
    imageSearch.setResultSetSize(ResultSetSize.LARGE);
    
    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
      public void onSearchResults(SearchResultsEvent resultsEvent) {
        handleSearchResults(event.query, resultsEvent.getResults());
      }
    });
    
    imageSearch.execute(event.query);
    currentQuery = event.query;
  }
  
  public void handleSearchResults(String query, JsArray<? extends Result> results) {
    if (!query.equals(currentQuery)) {
      return;
    }
    
    List<String> urls = new ArrayList<String>();
    
    for (int i = 0; i < results.length(); i++) {
      ImageResult r = (ImageResult) results.get(i);
      urls.add(r.getThumbnailUrl());
    }

    eventBus.fireEvent(new ImageResultsAvailableEvent(
        !query.equals(lastReceivedQuery), query, urls));    

    lastReceivedQuery = query;
  }
}
