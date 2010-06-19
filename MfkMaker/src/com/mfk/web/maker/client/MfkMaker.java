package com.mfk.web.maker.client;

import java.util.Vector;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.core.client.JsArray;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyCodes;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.search.client.ExpandMode;
import com.google.gwt.search.client.ImageResult;
import com.google.gwt.search.client.ImageSearch;
import com.google.gwt.search.client.Result;
import com.google.gwt.search.client.ResultSetSize;
import com.google.gwt.search.client.SafeSearchValue;
import com.google.gwt.search.client.SearchControl;
import com.google.gwt.search.client.SearchControlOptions;
import com.google.gwt.search.client.SearchResultsHandler;
import com.google.gwt.search.client.SearchStartingHandler;
import com.google.gwt.search.client.WebSearch;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.UIObject;
import com.google.gwt.user.client.ui.VerticalPanel;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkMaker implements EntryPoint {
	/**
	 * This is the entry point method.
	 */
	
	public static Vector<Image> results = new Vector<Image>();
	public static Image selected;
	
	public void onModuleLoad() {
		final RootPanel resultPanel = RootPanel.get("search-results");
		
	    SearchControlOptions options = new SearchControlOptions();
	    ImageSearch imageSearch = new ImageSearch();
	    imageSearch.setSafeSearch(SafeSearchValue.STRICT);
	    imageSearch.setResultSetSize(ResultSetSize.LARGE);
	    options.add(imageSearch, ExpandMode.CLOSED);
	    final SearchControl control = new SearchControl(options);
	    
	    final ResultClickHandler resultClick = new ResultClickHandler();
	    
	    control.addSearchStartingHandler(new SearchStartingHandler() {
			@Override
			public void onSearchStarting(SearchStartingEvent event) {
				resultPanel.clear();
				MfkMaker.results.clear();
			}
	    });
	    control.addSearchResultsHandler(new SearchResultsHandler() {
			@Override
			public void onSearchResults(SearchResultsEvent event) {
				JsArray<? extends Result> results = event.getResults();
				System.out.println("Handler! #results = " + results.length());
				for (int i = 0; i < results.length(); i++) {
					ImageResult r = (ImageResult)results.get(i);
					Image thumb = new Image(r.getThumbnailUrl());
					thumb.setHeight(String.valueOf(r.getThumbnailHeight()));
					thumb.setWidth(String.valueOf(r.getThumbnailWidth()));
					thumb.addStyleName("search-result");
					thumb.addClickHandler(resultClick);
					resultPanel.add(thumb);
					MfkMaker.results.add(thumb);
				}
			}
	    });
	    control.execute("treehouse");
	    RootPanel.get("search-control").add(control);
	    
	    final Button set1 = new Button("Set Item 1");
	    final Button set2 = new Button("Set Item 2");
	    final Button set3 = new Button("Set Item 3");
	    
	    set1.setEnabled(false);
	    set2.setEnabled(false);
	    set3.setEnabled(false);
	    
	    RootPanel.get("saved-1").add(set1);
	    RootPanel.get("saved-2").add(set2);
	    RootPanel.get("saved-3").add(set3);
	}
}

class ResultClickHandler implements ClickHandler {
	@Override
	public void onClick(ClickEvent event) {
		for (Image img: MfkMaker.results) {
			img.setStylePrimaryName("search-result");
		}
		
		Image source = (Image)event.getSource();
		source.setStylePrimaryName("search-result-sel");
		MfkMaker.selected = source;

		System.out.println("Clicked on " + source.getUrl());
	}
}


