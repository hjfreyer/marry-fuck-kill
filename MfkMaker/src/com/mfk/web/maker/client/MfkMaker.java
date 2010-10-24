package com.mfk.web.maker.client;

import java.util.Vector;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.core.client.JsArray;
import com.google.gwt.event.dom.client.ChangeEvent;
import com.google.gwt.event.dom.client.ChangeHandler;
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
import com.google.gwt.user.client.Timer;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.HTMLPanel;
import com.google.gwt.user.client.ui.HorizontalPanel;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.Panel;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.UIObject;
import com.google.gwt.user.client.ui.VerticalPanel;
import com.google.gwt.user.client.ui.Widget;

// TODO(mjkelly): See this example:
// https://code.google.com/apis/ajax/playground/#raw_search

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkMaker implements EntryPoint {
	/**
	 * This is the entry point method.
	 */
	
	public static Image selected;

	// arrays holding the attributes of the 3 items in the new triple
	public static TextBox names[] = {new TextBox(),
		                             new TextBox(),
		                             new TextBox()};
	public static Image images[] = {null, null, null};
	public static Button setButtons[] = {new Button("Set Item 1"),
									     new Button("Set Item 2"),
									     new Button("Set Item 3")};
	
	public static MfkPanel items[] = {null, null, null};
	
	// the actual image results
	public static Vector<Image> results = new Vector<Image>();
	// the UI panel that displays the search results
    static Panel resultPanel;
    
    // a pane shown only when loading results
    static RootPanel resultsLoadingPanel = RootPanel.get("results-loading");
    
	static ImageSearch imageSearch = new ImageSearch();
	
	static public Button searchButton = new Button("Search");
	static public TextBox searchBox = new TextBox();
	
	static final EditDialog editDialog = new EditDialog(true);
	
	static final HorizontalPanel itemPanel = new HorizontalPanel();
	
	static final String DEFAULT_SEARCH = "treehouse";
	static final HTML LOADING =
		new HTML("<img src=\"/gwt/loading.gif\" alt=\"\"> Loading...");

	public void onModuleLoad() {
		MfkMaker.resultsLoadingPanel.setVisible(false);
		MfkMaker.searchButton = new Button("Search");
		MfkMaker.searchBox = new TextBox();
		MfkMaker.resultPanel = new FlowPanel();
		MfkMaker.imageSearch = new ImageSearch();
		RootPanel.get("created-items").add(MfkMaker.itemPanel);
		MfkMaker.itemPanel.setSpacing(10);
		
		MfkMaker.names[0].setText("item one name");
		MfkMaker.names[1].setText("item two name");
		MfkMaker.names[2].setText("item three name");
		
		for (int i = 1; i <= 3; i++) {
			Image img = new Image("/gwt/images/treehouse-" + i + ".jpeg");
			MfkPanel item = new MfkPanel("Treehouse " + i, img);
			MfkMaker.addItem(item);
		}
		
	    final SearchControlOptions options = new SearchControlOptions();
	    
	    imageSearch.setSafeSearch(SafeSearchValue.STRICT);
	    imageSearch.setResultSetSize(ResultSetSize.LARGE);
	    options.add(imageSearch, ExpandMode.OPEN);
	    options.setKeepLabel("<b>Keep It!</b>");
	    options.setLinkTarget(LinkTarget.BLANK);
	    MfkMaker.editDialog.setAnimationEnabled(true);
	    
	    final ClickHandler resultClick = new ClickHandler() {
			public void onClick(ClickEvent event) {
				MfkMaker.editDialog.setImage((Image)event.getSource());
			}
	    };
	    
	    // This handles the displayed result list.
	    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
			public void onSearchResults(SearchResultsEvent event) {
				MfkMaker.resultsLoadingPanel.setVisible(false);
				JsArray<? extends Result> results = event.getResults();
				System.out.println("List handler! #results = " + results.length());
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
	    
	    // This handles the auto-set image.
	    imageSearch.addSearchResultsHandler(new SearchResultsHandler() {
			public void onSearchResults(SearchResultsEvent event) {
				JsArray<? extends Result> results = event.getResults();
				System.out.println("Top-result handler! #results = " + results.length());
				if (results.length() >= 1) {
					ImageResult r = (ImageResult)results.get(0);
					Image image = new Image(r.getThumbnailUrl());
					MfkMaker.editDialog.setImage(image);
				}
			}
	    });
	    
	    searchButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				MfkMaker.doSearch();
			}
	    });
	    searchBox.addKeyPressHandler(new KeyPressHandler() {
			public void onKeyPress(KeyPressEvent event) {
				if (event.getCharCode() == '\n' || event.getCharCode() == '\r') {
					MfkMaker.doSearch();
				}
			}
	    });
	}
	
	public static void doSearch() {
		MfkMaker.resultsLoadingPanel.setVisible(true);
		MfkMaker.resultPanel.clear();
		MfkMaker.results.clear();
		MfkMaker.imageSearch.execute(MfkMaker.searchBox.getText());
	}
	
	/**
	 * Add an item to the page.
	 * @param item the MfkPanel to add
	 */
	public static void addItem(MfkPanel item) {
		item.addToPanel(MfkMaker.itemPanel);
	}
	
	/**
	 * Update the page's status with instructions, or how many items left to
	 * create.
	 */
	public static void updateStatus() {
		if (MfkPanel.count > 0) {
			MfkMaker.setStatus((3 - MfkPanel.count) + " items remaining.");
		}
		else {
			MfkMaker.setStatus("Click an image to create an item.");
		}
	}
	
	/**
	 * Set the page's status.
	 * @param s status string
	 */
	private static void setStatus(String s) {
		Panel counter = RootPanel.get("counter");
		counter.clear();
		counter.add(new HTML("<h2>" + s + "</h2>"));
	}
}

class EditDialog extends DialogBox {
	private MfkPanel item = null;
	private Image editImage = new Image();
	private TextBox editTitle = new TextBox();
	
	// These are all bookkeeping for auto-search.
	private long lastSearchTimeMillis = 0;
	private String lastSearch = "";
	private Timer timer;
	
	public EditDialog(boolean b) {
		super(b);
	}

	public void editItem(final MfkPanel item) {
		// TODO Auto-generated method stub
		this.item = item;
		System.out.println("Showing dialog for :" + item);
		this.editImage.setUrl(item.image.getUrl());
		this.editTitle.setText(item.title);
		
		// TODO: We are *not* guaranteed to get the final state until the
		// search box loses focus! Fix this if deploying this UI!
		this.editTitle.addChangeHandler(new ChangeHandler() {
			public void onChange(ChangeEvent event) {
				maybeSetImageFromTitle();
			}
		});
		this.editTitle.addKeyPressHandler(new KeyPressHandler() {
			public void onKeyPress(KeyPressEvent event) {
				maybeSetImageFromTitle();
			}
		});
		final Panel search = new VerticalPanel();
		VerticalPanel p = new VerticalPanel();
		p.setSpacing(5);
		// TODO: put this in CSS, come up with a well-reasoned value
		p.setWidth("600px");
		
		Button create = new Button("<b>Save</b>");
		create.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent e) {
				System.out.println("Should create item here");
				hide();
				// update the existing item
				item.setTitle(editTitle.getText());
				item.setImage(editImage);
			}
		});
		
		Button cancel = new Button("Cancel");
		cancel.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent e) {
				hide();
			}
		});
		
		p.add(new HTML("<b>Name:</b>"));
		p.add(editTitle);
		p.add(new HTML("<b>Image:</b>"));
		p.add(editImage);
		p.add(new HTML("Not the image you wanted?"));
		HTML link = new HTML("See more images.");
		link.addClickHandler(new ClickHandler (){
			public void onClick(ClickEvent event) {
				search.setVisible(!search.isVisible());
				if (search.isVisible()) {
					MfkMaker.searchBox.setText(editTitle.getText());
					MfkMaker.doSearch();
				}
			}
		});
		link.setStylePrimaryName("fakelink");
		p.add(link);
		
		HorizontalPanel searchControls = new HorizontalPanel();
		searchControls.add(MfkMaker.searchBox);
		searchControls.add(MfkMaker.searchButton);
		search.add(new HTML("<hr>Search for more images:"));
		search.add(searchControls);
		search.add(MfkMaker.resultPanel);
		search.setVisible(false);
		
		HorizontalPanel buttonPanel = new HorizontalPanel();
		buttonPanel.setSpacing(5);
		buttonPanel.add(create);
		buttonPanel.add(cancel);
		p.add(buttonPanel);
		p.add(search);
		
		this.setWidget(p);
		this.show();
		this.center();
		this.editTitle.setFocus(true);
	}
	
	/**
	 * Check the time since the last time we auto-searched, perform a
	 * search for the title string, and set the image to the top result.
	 * 
	 * If don't perform the search if we just searched recently.
	 */
	public void maybeSetImageFromTitle() {
		long now = System.currentTimeMillis();
		String text = this.editTitle.getText();
		
		if (now - this.lastSearchTimeMillis > 1000) {
			this.lastSearchTimeMillis = now;
			if (!text.equals(this.lastSearch) && !text.isEmpty()) {
				this.lastSearch = text;
				MfkMaker.searchBox.setText(text);
				MfkMaker.doSearch();
				System.out.println("Auto-search: " + text + ">");
			}
		}
		
		if (this.timer == null && !text.equals(this.lastSearch)) {
			this.timer = new Timer() {
				public void run() {
					System.out.println("Timer!");
					timer = null;
					maybeSetImageFromTitle();
				}
			};
			this.timer.schedule(1200);
		}
	}
	
	/**
	 * Set the image for the item under edit.
	 */
	public void setImage(Image image) {
		System.out.println("Set edit image url = " + image.getUrl());
		this.editImage.setUrl(image.getUrl());
	}
	
	public void hide() {
		super.hide();
		this.item = null;
	}
	
	public MfkPanel getItem() {
		return this.item;
	}
}

class MfkPanel extends VerticalPanel {
	public String title;
	public Image image = new Image();
	private Panel parent;
	
	// How many MfkPanels have been shown (i.e., added to another panel).
	static int count = 0;
	
	public MfkPanel(String title, Image image) {
		this.setTitle(title);
		this.setImage(image);
		System.out.println("MfkPanel: title:" + title +
				           ", count:" + MfkPanel.count); 
	}
	
	public void setImage(Image image) {
		this.image.setUrl(image.getUrl());
		this.refresh();
	}
	
	public void setTitle(String title) {
		this.title = title;
		this.refresh();
	}
	
	/**
	 * Add this object to another panel. This is a grab-bag of misc logic
	 * associated with the MfkMaker.
	 * @param p
	 */
	public void addToPanel(Panel p) {
		this.parent = p;
		MfkPanel.count++;
		this.parent.add(this);
		MfkMaker.updateStatus();
	}
	
	/**
	 * Remove this object from whatever panel it was added to. (It *must* have
	 * been previously added.)
	 */
	public void remove() {
		this.parent.remove(this);
		MfkPanel.count--;
		MfkMaker.updateStatus();
	}
	
	/**
	 * Refresh the UI elements of the page.
	 */
	private void refresh() {
		this.clear();
		final MfkPanel outerThis = this;
		Button editButton = new Button("Edit");
		editButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				System.out.println("Delete " + this);
				MfkMaker.editDialog.editItem(outerThis);
			}
		});
		this.add(editButton);
		this.add(new HTML("<i>" + this.title + "</i>"));
		this.add(this.image);
	}
	
	public String toString() {
		return "<MfkPanel: " + this.title +
		       ", url=" + this.image.getUrl() + ">";
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