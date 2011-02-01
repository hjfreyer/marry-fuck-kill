package com.mfk.web.maker.client;

import java.util.Vector;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.JsArray;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.search.client.ExpandMode;
import com.google.gwt.search.client.ImageResult;
import com.google.gwt.search.client.ImageSearch;
import com.google.gwt.search.client.LinkTarget;
import com.google.gwt.search.client.Result;
import com.google.gwt.search.client.ResultSetSize;
import com.google.gwt.search.client.SafeSearchValue;
import com.google.gwt.search.client.SearchControlOptions;
import com.google.gwt.search.client.SearchResultsHandler;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.HorizontalPanel;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Panel;
import com.google.gwt.user.client.ui.RootPanel;
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

  // the actual image results
  public static Vector<Image> results = new Vector<Image>();
  // the UI panel that displays the search results
  static Panel resultPanel = new FlowPanel();
  // The query that led to the displayed results.
  public static String resultsSearchQuery;

  static ImageSearch imageSearch = new ImageSearch();

  static final EditDialog editDialog = new EditDialog(true);
  static Vector<MfkPanel> items = new Vector<MfkPanel>();

  static final HorizontalPanel itemPanel = new HorizontalPanel();

  static final DialogBox miscDialog = new DialogBox();

  public static final String[] defaultURLs = {
    "http://www.marryboffkill.net/s/treehouse-1.jpeg",
    "http://www.marryboffkill.net/s/treehouse-2.jpeg",
    "http://www.marryboffkill.net/s/treehouse-3.jpeg"};

  public void onModuleLoad() {
    RootPanel.get("created-items").add(MfkMaker.itemPanel);
    MfkMaker.itemPanel.setSpacing(10);

    for (int i = 0; i < 3; i++) {
      SearchImage img = new SearchImage(defaultURLs[i], "treehouse");
      MfkPanel item = new MfkPanel("Treehouse " + (i+1), img);
      MfkMaker.addItem(item);
    }

    final SearchControlOptions options = new SearchControlOptions();

    // TODO(mjkelly): reconsider this value. Remember to synchronize any
    // change with the server-side checking.
    imageSearch.setSafeSearch(SafeSearchValue.MODERATE);
    imageSearch.setResultSetSize(ResultSetSize.LARGE);
    options.add(imageSearch, ExpandMode.OPEN);
    options.setKeepLabel("<b>Keep It!</b>");
    options.setLinkTarget(LinkTarget.BLANK);
    MfkMaker.editDialog.setAnimationEnabled(true);

    final ClickHandler resultClick = new ClickHandler() {
      public void onClick(ClickEvent event) {
        Image source = (Image) event.getSource();
        MfkMaker.editDialog.setImage(source);
      }
    };

    // This handles the displayed result list.
    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
      public void onSearchResults(SearchResultsEvent event) {
        JsArray<? extends Result> results = event.getResults();
        System.out.println("List handler! #results = "
            + results.length());
        for (int i = 0; i < results.length(); i++) {
          ImageResult r = (ImageResult) results.get(i);
          Image thumb = new Image(r.getThumbnailUrl());
          thumb.setHeight(String.valueOf(r.getThumbnailHeight()));
          thumb.setWidth(String.valueOf(r.getThumbnailWidth()));
          thumb.addStyleName("search-result");
          thumb.addClickHandler(resultClick);
          resultPanel.add(thumb);
          MfkMaker.results.add(thumb);
          MfkMaker.editDialog.setSearchThrobber(false);
        }
      }
    });

    // This handles the auto-set image.
    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
      public void onSearchResults(SearchResultsEvent event) {
        JsArray<? extends Result> results = event.getResults();
        System.out.println("Top-result handler! #results = "
            + results.length() + ", search = "
            + MfkMaker.resultsSearchQuery);
        if (results.length() >= 1) {
          ImageResult r = (ImageResult) results.get(0);
          Image image = new Image(r.getThumbnailUrl());
          MfkMaker.editDialog.autoSetImage(image);
          MfkMaker.editDialog.setAutoThrobber(false);
        }
      }
    });

    // The submit button.
    Button submit_btn = new Button("Create");
    submit_btn.addClickHandler(new SubmitHandler());
    RootPanel.get("submit").add(submit_btn);
  }

  /**
   * Add an item to the page.
   *
   * @param item
   *      the MfkPanel to add
   */
  public static void addItem(MfkPanel item) {
    MfkMaker.items.add(item);
    MfkMaker.itemPanel.add(item);
  }

  public static void showDialog(String title, String message) {
    MfkMaker.miscDialog.hide();
    MfkMaker.miscDialog.clear();
    VerticalPanel panel = new VerticalPanel();
    MfkMaker.miscDialog.add(panel);

    panel.add(new HTML("<p><b>" + title + "</b></p>"));
    panel.add(new HTML("<p>" + message + "</p>"));

    Button closeButton = new Button("OK");
    closeButton.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent event) {
        MfkMaker.miscDialog.hide();
      }
    });
    panel.add(closeButton);
    MfkMaker.miscDialog.center();
    MfkMaker.miscDialog.show();
    closeButton.setFocus(true);
  }
}
