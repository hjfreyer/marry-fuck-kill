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
import com.google.gwt.user.client.ui.Panel;
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
	
	public static TextBox names[] = {new TextBox(),
		                             new TextBox(),
		                             new TextBox()};
	
	public static Image images[] = {null, null, null};
	
	public static Button setButtons[] = {new Button("Set Item 1"),
									     new Button("Set Item 2"),
									     new Button("Set Item 3")};
	
	public void onModuleLoad() {
		final RootPanel resultPanel = RootPanel.get("search-results");
		
		
		MfkMaker.names[0].setText("item one name");
		MfkMaker.names[1].setText("item two name");
		MfkMaker.names[2].setText("item three name");
		
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
	    
	    for (int i = 0; i < 3; i++) {
	    	MfkMaker.setButtons[i].setEnabled(false);
	    	MfkMaker.setButtons[i].addClickHandler(new SetImageHandler(i));
	    }
	    
	    RootPanel.get("saved-1-name").add(names[0]);
	    RootPanel.get("saved-1-button").add(MfkMaker.setButtons[0]);
	    
	    RootPanel.get("saved-2-name").add(names[1]);
	    RootPanel.get("saved-2-button").add(MfkMaker.setButtons[1]);
	    
	    RootPanel.get("saved-3-name").add(names[2]);
	    RootPanel.get("saved-3-button").add(MfkMaker.setButtons[2]);
	    
	    
	    final Button submitButton = new Button("Create!");
	    submitButton.addClickHandler(new SubmitHandler());
	    
	    RootPanel.get("submit-button").add(submitButton);
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

		for (Button b: MfkMaker.setButtons) {
			b.setEnabled(true);
		}
		
		System.out.println("Clicked on " + source.getUrl());
	}
}

class SetImageHandler implements ClickHandler {
	private int itemIndex;
	private String id;
	public SetImageHandler(int itemIndex) {
		this.itemIndex = itemIndex;
		this.id = "saved-" + Integer.toString(this.itemIndex+1) + "-inner";
	}

	@Override
	public void onClick(ClickEvent event) {
		System.out.println("Click #" + this.itemIndex + " -> " + this.id);
		RootPanel p = RootPanel.get(this.id);
		p.clear();
		
		Image img = new Image(MfkMaker.selected.getUrl());
		p.add(img);
		MfkMaker.images[this.itemIndex] = img;
	}
}

class SubmitHandler implements ClickHandler {
	@Override
	public void onClick(ClickEvent event) {
		System.out.println("Want to create:\n"
				+ "{n:" + MfkMaker.names[0].getText()
				+ ", u:" + MfkMaker.images[0].getUrl() + "}\n"
				+ "{n:" + MfkMaker.names[1].getText()
				+ ", u:" + MfkMaker.images[1].getUrl() + "}\n"
				+ "{n:" + MfkMaker.names[2].getText()
				+ ", u:" + MfkMaker.images[2].getUrl() + "}");
	}
}