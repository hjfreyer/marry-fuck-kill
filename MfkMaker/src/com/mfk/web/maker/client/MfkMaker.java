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

	// XXX: These aren't used in the current UI.
	// TODO: Remove them and rewrite the RPC logic to use MfkPanels.
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
    static Panel resultPanel = new FlowPanel();
    
	static ImageSearch imageSearch = new ImageSearch();
	
	static final EditDialog editDialog = new EditDialog(true);
	
	static final HorizontalPanel itemPanel = new HorizontalPanel();
	
	public void onModuleLoad() {
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
					MfkMaker.editDialog.setSearchThrobber(false);
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
					MfkMaker.editDialog.autoSetImage(image);
					MfkMaker.editDialog.setAutoThrobber(false);
				}
			}
	    });
	}
	

	/**
	 * Add an item to the page.
	 * @param item the MfkPanel to add
	 */
	public static void addItem(MfkPanel item) {
		MfkMaker.itemPanel.add(item);
	}
}

class EditDialog extends DialogBox {
	private static String THROBBER_URL = "/gwt/loading.gif";
	private MfkPanel item = null;
	private Image editImage = new Image();
	private TextBox editTitle = new TextBox();
	private Image autoThrobber = new Image(EditDialog.THROBBER_URL);
	private HorizontalPanel searchThrobber = new HorizontalPanel();
	
	private Button searchButton = new Button("Search");
	private TextBox searchBox = new TextBox();
	
	// These are all bookkeeping for auto-search:
	
	// Last time we sent a search. Maintained by maybeSearch.
	private long lastSearchMillis = 0;
	// Last time the text box changed. Maintained by repeatingTimer.
	private long lastChangeMillis = 0;
	// Last search text. Maintained by maybeSearch.
	private String lastSearch = "";
	
	// Timer that drives actual searching.
	private Timer repeatingTimer;
	
	// The expanding search panel.
	private VerticalPanel search = new VerticalPanel();
	
	public EditDialog() {
		this(false);
	}
	
	public EditDialog(boolean b) {
		super(b);
		
		HorizontalPanel searchControls = new HorizontalPanel();
		
		this.searchButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				doSearch();
			}
	    });
	    this.searchBox.addKeyPressHandler(new KeyPressHandler() {
			public void onKeyPress(KeyPressEvent event) {
				if (event.getCharCode() == '\n' || event.getCharCode() == '\r') {
					doSearch();
				}
			}
	    });
		
		searchControls.add(this.searchBox);
		searchControls.add(this.searchButton);
		
		HorizontalPanel moreImagesTitle = new HorizontalPanel();
		moreImagesTitle.add(new HTML("Search for more images:"));
		HTML hideLink = new HTML("[hide]");
		hideLink.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				search.setVisible(false);
			}
		});
		hideLink.setStylePrimaryName("fakelink");
		moreImagesTitle.add(hideLink);
		
		this.search.add(new HTML("<hr>"));
		this.search.add(moreImagesTitle);
		this.search.add(searchControls);
		
		this.searchThrobber.add(new Image(EditDialog.THROBBER_URL));
		this.searchThrobber.add(new HTML("Loading..."));
		this.searchThrobber.setSpacing(10);
		this.searchThrobber.setVisible(false);
		this.search.add(searchThrobber);
		
		this.search.add(MfkMaker.resultPanel);
		this.search.setSpacing(5);
	}

	public void editItem(final MfkPanel item) {
		this.item = item;
		System.out.println("Showing dialog for :" + item);
		this.editImage.setUrl(item.image.getUrl());
		this.editTitle.setText(item.title);
		this.autoThrobber.setVisible(false);
		this.search.setVisible(false);
		
		long now = System.currentTimeMillis();
		this.lastSearch = item.title;
		this.lastChangeMillis = this.lastSearchMillis = now;

		// This just keeps track of when the last change in the box was.
		// If it misses a keystroke, our time is a little old, but that's okay.
		// (We still throttle searches.)
		this.editTitle.addKeyPressHandler(new KeyPressHandler() {
			public void onKeyPress(KeyPressEvent event) {
				lastChangeMillis = System.currentTimeMillis();
				setAutoThrobber(true);
			}
		});

		// This checks repeatedly if we should search.
		this.repeatingTimer = new Timer() {
			public void run() {
				maybeSearch();
			}
		};
		this.repeatingTimer.scheduleRepeating(250);
		
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
		HorizontalPanel titlePanel = new HorizontalPanel();
		titlePanel.add(editTitle);
		titlePanel.add(autoThrobber);
		p.add(titlePanel);
		p.add(new HTML("<b>Image:</b>"));
		p.add(editImage);
		p.add(new HTML("Not the image you wanted?"));
		HTML link = new HTML("See more images.");
		link.addClickHandler(new ClickHandler (){
			public void onClick(ClickEvent event) {
				search.setVisible(!search.isVisible());
				if (search.isVisible()) {
					searchBox.setText(editTitle.getText());
					doSearch();
				}
			}
		});
		link.setStylePrimaryName("fakelink");
		p.add(link);
		
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
	 * This is the logic that determines if we should search.
	 */
	public void maybeSearch() {
		long now = System.currentTimeMillis();
		String text = this.editTitle.getText();
		
		if (now - this.lastChangeMillis > 250) {
			if (now - this.lastSearchMillis > 1000) {
				this.lastSearchMillis = now;
				if (!text.equals(this.lastSearch) && !text.isEmpty()) {
					this.lastSearch = text;
					System.out.println("maybeSearch: <" + text + ">");
					this.doAutoSearch();
				}
			}
		}
	}
	
	private void doAutoSearch() {
		String text = this.editTitle.getText();
		this.searchBox.setText(text);
		this.doSearch();
	}
	
	private void doSearch() {
		MfkMaker.editDialog.setSearchThrobber(true);
		MfkMaker.resultPanel.clear();
		MfkMaker.results.clear();
		MfkMaker.imageSearch.execute(this.searchBox.getText());
	}
	
	/**
	 * Turn on or off the throbber.
	 * @param enabled
	 */
	public void setAutoThrobber(boolean enabled) {
		this.autoThrobber.setVisible(enabled);
	}
	
	public void setSearchThrobber(boolean enabled) {
		this.searchThrobber.setVisible(enabled);
	}
	
	/**
	 * Set the image for the item under edit.
	 */
	public void setImage(Image image) {
		System.out.println("Set edit image url = " + image.getUrl());
		this.editImage.setUrl(image.getUrl());
	}
	
	public void autoSetImage(Image image) {
		if (this.shouldAutoSet()) {
			this.setImage(image);
		}
	}
	
	private boolean shouldAutoSet() {
		return !this.search.isVisible();
	}
	
	public void hide() {
		super.hide();
		this.item = null;
		if (this.repeatingTimer != null)
			this.repeatingTimer.cancel();
	}
	
	public MfkPanel getItem() {
		return this.item;
	}
}

class MfkPanel extends VerticalPanel {
	public String title;
	public Image image = new Image();
	private ClickHandler editMe;
	
	public MfkPanel(String title, Image image) {
		this.setTitle(title);
		this.setImage(image);
		System.out.println("MfkPanel: title:" + title); 
		this.addStyleName("mfkpanel");
		
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
	 * Refresh the UI elements of the page.
	 */
	private void refresh() {
		this.clear();
		Button editButton = new Button("Edit");
		final MfkPanel outerThis = this;
		editButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				MfkMaker.editDialog.editItem(outerThis);
			}
		});
		this.add(new HTML("<i>" + this.title + "</i>"));
		this.add(this.image);
		this.add(editButton);
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

