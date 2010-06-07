package com.mfk.web.maker.client;

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
import com.google.gwt.search.client.SearchControl;
import com.google.gwt.search.client.SearchControlOptions;
import com.google.gwt.search.client.WebSearch;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.VerticalPanel;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkMaker implements EntryPoint {
	/**
	 * This is the entry point method.
	 */
	public void onModuleLoad() {
		Button b = new Button("Foo!");
		RootPanel.get("main").add(b);
		System.out.println("MfkMaker start");
		
		// this is how  I think part of it might look for us:
//		ImageSearch search = new ImageSearch();
//		search.execute("Google");
//		
//		JsArray<Result> r = (JsArray<Result>) search.getResults();
//		
//		System.out.println(r.length() + " results for query.");
//		for (int i = 0 ; i < r.length(); i++) {
//			Result img = r.get(i);
//			System.out.println("Result #" + i + ": " + img);
//		}
		
	    SearchControlOptions options = new SearchControlOptions();
	    WebSearch webSearch = new WebSearch();
	    webSearch.setResultSetSize(ResultSetSize.LARGE);
	    options.add(webSearch);
	    ImageSearch imageSearch = new ImageSearch();
	    options.add(imageSearch, ExpandMode.OPEN);
	    final SearchControl control = new SearchControl(options);
	    control.execute("treehouse");
	    RootPanel.get().add(control);
		
		System.out.println("MfkMaker end");
	}
}
