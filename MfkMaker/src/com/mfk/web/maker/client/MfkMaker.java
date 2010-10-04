package com.mfk.web.maker.client;

import java.util.Vector;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.core.client.JsArray;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyCodes;
import com.google.gwt.event.dom.client.KeyPressEvent;
import com.google.gwt.event.dom.client.KeyPressHandler;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.http.client.Request;
import com.google.gwt.http.client.RequestBuilder;
import com.google.gwt.http.client.RequestCallback;
import com.google.gwt.http.client.RequestException;
import com.google.gwt.http.client.Response;
import com.google.gwt.search.client.DrawMode;
import com.google.gwt.search.client.ExpandMode;
import com.google.gwt.search.client.ImageResult;
import com.google.gwt.search.client.ImageSearch;
import com.google.gwt.search.client.KeepHandler;
import com.google.gwt.search.client.LinkTarget;
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

// TODO(mjkelly): See this example:
// https://code.google.com/apis/ajax/playground/#raw_search

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
	
	final static Button searchButton = new Button("Search");
    final static TextBox searchBox = new TextBox();
    final static RootPanel resultPanel = RootPanel.get("search-results");
	final static ImageSearch imageSearch = new ImageSearch();
	
	public void onModuleLoad() {
		MfkMaker.names[0].setText("item one name");
		MfkMaker.names[1].setText("item two name");
		MfkMaker.names[2].setText("item three name");
		
	    final SearchControlOptions options = new SearchControlOptions();
	    imageSearch.setSafeSearch(SafeSearchValue.STRICT);
	    imageSearch.setResultSetSize(ResultSetSize.LARGE);
	    options.add(imageSearch, ExpandMode.OPEN);
	    options.setKeepLabel("<b>Keep It!</b>");
	    options.setLinkTarget(LinkTarget.BLANK);
	    final ResultClickHandler resultClick = new ResultClickHandler();
	    
	    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
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
	    imageSearch.execute("treehouse");
	    
	    searchBox.setText("treehouse");
	    
	    searchButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				resultPanel.clear();
				MfkMaker.results.clear();
				imageSearch.execute(searchBox.getText());
			}
	    });
	    
	    searchBox.addKeyPressHandler(new KeyPressHandler() {
			public void onKeyPress(KeyPressEvent event) {
				if (event.getCharCode() == '\n' || event.getCharCode() == '\r') {
					resultPanel.clear();
					MfkMaker.results.clear();
					imageSearch.execute(searchBox.getText());
				}
			}
	    });
	    
	    RootPanel.get("search-control").add(searchBox);
	    RootPanel.get("search-control").add(searchButton);
	    
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

	public void onClick(ClickEvent event) {
		System.out.println("Click #" + this.itemIndex + " -> " + this.id);
		RootPanel p = RootPanel.get(this.id);
		p.clear();
		
		Image img = new Image(MfkMaker.selected.getUrl());
		p.add(img);
		MfkMaker.images[this.itemIndex] = img;
	}
}

// TODO(mjkelly): Do client-side validation here.
class SubmitHandler implements ClickHandler {
	public void onClick(ClickEvent event) {
		System.out.println("Want to create:\n"
				+ "{n:" + MfkMaker.names[0].getText()
				+ ", u:" + MfkMaker.images[0].getUrl() + "}\n"
				+ "{n:" + MfkMaker.names[1].getText()
				+ ", u:" + MfkMaker.images[1].getUrl() + "}\n"
				+ "{n:" + MfkMaker.names[2].getText()
				+ ", u:" + MfkMaker.images[2].getUrl() + "}");
		
		
		String url = "/rpc/create/";
		RequestBuilder builder = new RequestBuilder(RequestBuilder.POST, url);
		builder.setHeader("Content-Type", "application/x-www-form-urlencoded");
		StringBuffer reqData = new StringBuffer();
		
		reqData.append("n1=").append(MfkMaker.names[0].getText());
		reqData.append("&n2=").append(MfkMaker.names[1].getText());
		reqData.append("&n3=").append(MfkMaker.names[2].getText());
		reqData.append("&u1=").append(MfkMaker.images[0].getUrl());
		reqData.append("&u2=").append(MfkMaker.images[1].getUrl());
		reqData.append("&u3=").append(MfkMaker.images[2].getUrl());
		
		try {
			builder.sendRequest(reqData.toString(), new RequestCallback() {
				public void onError(Request request, Throwable exception) {
					System.out.println("Error creating new Triple");
				}

				public void onResponseReceived(Request request,
						Response response) {
					if (response.getStatusCode() == 200) {
						System.out.println("Successful creation request: "
								+ response.getText());
					}
					else {
						System.out.println("Server didn't like our new triple. "
								+ "Response code: " + response.getStatusCode()
								+ ". Response text: " + response.getText());
					}
				}
			});
		} catch (RequestException e) {
			System.out.println("Error sending vote: " + e);
		}
		
	}
}