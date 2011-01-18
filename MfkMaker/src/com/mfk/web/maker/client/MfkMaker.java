package com.mfk.web.maker.client;

import java.util.Vector;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.JsArray;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyPressEvent;
import com.google.gwt.event.dom.client.KeyPressHandler;
import com.google.gwt.http.client.Request;
import com.google.gwt.http.client.RequestBuilder;
import com.google.gwt.http.client.RequestCallback;
import com.google.gwt.http.client.RequestException;
import com.google.gwt.http.client.Response;
import com.google.gwt.http.client.URL;
import com.google.gwt.http.client.UrlBuilder;
import com.google.gwt.search.client.ExpandMode;
import com.google.gwt.search.client.ImageResult;
import com.google.gwt.search.client.ImageSearch;
import com.google.gwt.search.client.LinkTarget;
import com.google.gwt.search.client.Result;
import com.google.gwt.search.client.ResultSetSize;
import com.google.gwt.search.client.SafeSearchValue;
import com.google.gwt.search.client.SearchControlOptions;
import com.google.gwt.search.client.SearchResultsHandler;
import com.google.gwt.user.client.Timer;
import com.google.gwt.user.client.Window;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.HorizontalPanel;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.Panel;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
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
		"http://www.marryboffkill.net/static/treehouse-1.jpeg",
		"http://www.marryboffkill.net/static/treehouse-2.jpeg",
		"http://www.marryboffkill.net/static/treehouse-3.jpeg"};

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
		imageSearch.setSafeSearch(SafeSearchValue.STRICT);
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
	 *            the MfkPanel to add
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

class EditDialog extends DialogBox {
	private static String THROBBER_URL = "/loading.gif";
	private MfkPanel item = null;
	private SearchImage editImage = new SearchImage();
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
	// NOTE: This is distinct from MfkMaker.resultSearchQuery -- this is used
	// only for maybeSearch's retry logic.
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
		this.editImage.setUrlAndQuery(item.image.getUrl(), item.image
				.getQuery());
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
		link.addClickHandler(new ClickHandler() {
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
		MfkMaker.resultsSearchQuery = this.searchBox.getText();
		MfkMaker.imageSearch.execute(this.searchBox.getText());
	}

	/**
	 * Turn on or off the throbber.
	 * 
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
		System.out.println("Set edit image url = " + image.getUrl()
				+ " (from = " + MfkMaker.resultsSearchQuery + ")");
		this.editImage.setUrlAndQuery(image.getUrl(),
				MfkMaker.resultsSearchQuery);
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
	// The user-visible title for the entity
	public String title = "";
	// The user-visible image for the entity.
	public SearchImage image = new SearchImage();

	public MfkPanel(String title, SearchImage image) {
		this.setTitle(title);
		this.setImage(image);
		System.out.println("MfkPanel: title:" + title);
		this.addStyleName("mfkpanel");

	}

	public void setImage(SearchImage image) {
		this.image.setUrlAndQuery(image.getUrl(), image.getQuery());
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
		HTML title = new HTML(this.title);
		title.addStyleName("itemtitle");
		this.add(title);
		this.add(this.image);
		this.add(editButton);
	}

	public String toString() {
		return "<MfkPanel: " + this.title + ", url=" + this.image.getUrl()
				+ ">";
	}
	
	public boolean equals(MfkPanel o) {
		if (o == null)
			return false;
		return this.title.equals(o.title) &&
				this.image.getUrl().equals(o.image.getUrl());
	}
}

// TODO(mjkelly): Do client-side validation here.
class SubmitHandler implements ClickHandler {
	public void onClick(ClickEvent event) {
		MfkPanel[] p = { MfkMaker.items.get(0), MfkMaker.items.get(1),
				MfkMaker.items.get(2) };
		
		// Do some simple validation. This just to prevent mistakes. It
		// obviously does not replace real validation on the server side.
		if (p[0].equals(p[1]) || p[1].equals(p[2]) || p[2].equals(p[0])) {
			MfkMaker.showDialog("Error Creating MFK",
					"No two items can be identical. Please change one.");
			return;
		}
		for (int i = 0; i < 3; i++) {
			System.out.println("Checking URLs: " + p[i].image.getUrl() + " vs " + MfkMaker.defaultURLs[i]);
			if (p[i].image.getUrl().equals(MfkMaker.defaultURLs[i])) {
				MfkMaker.showDialog("Error Creating MFK",
						"You must change all the items from their defaults.");
				return;
			}
		}

		String url = "/rpc/create/";
		RequestBuilder builder = new RequestBuilder(RequestBuilder.POST, url);
		builder.setHeader("Content-Type", "application/x-www-form-urlencoded");
		StringBuffer reqData = new StringBuffer();

		URL.encodeQueryString(url);
		reqData.append("n1=").append(URL.encodeQueryString(p[0].title));
		reqData.append("&n2=").append(URL.encodeQueryString(p[1].title));
		reqData.append("&n3=").append(URL.encodeQueryString(p[2].title));
		reqData.append("&u1=").append(
				URL.encodeQueryString(p[0].image.getUrl()));
		reqData.append("&u2=").append(
				URL.encodeQueryString(p[1].image.getUrl()));
		reqData.append("&u3=").append(
				URL.encodeQueryString(p[2].image.getUrl()));
		reqData.append("&q1=").append(
				URL.encodeQueryString(p[0].image.getQuery()));
		reqData.append("&q2=").append(
				URL.encodeQueryString(p[1].image.getQuery()));
		reqData.append("&q3=").append(
				URL.encodeQueryString(p[2].image.getQuery()));
		System.out.println("request data = " + reqData);

		try {
			builder.sendRequest(reqData.toString(), new RequestCallback() {
				public void onError(Request request, Throwable exception) {
					System.out.println("Error creating new Triple");
				}

				public void onResponseReceived(Request request,
						Response response) {
					if (response.getStatusCode() == 200) {
						String[] responseParts = response.getText().split(":",
								2);
						if (responseParts[0].equals("ok")) {
							System.out.println("Successful creation request: "
									+ response.getText());
							
							// clear the existing page in a haphazard way
							MfkMaker.itemPanel.clear();
							RootPanel.get("submit").clear();
							RootPanel.get("instructions").setVisible(false);
							
							UrlBuilder builder = Window.Location.createUrlBuilder();
							builder.setPath("/");
							builder.setParameter("id", responseParts[1]);
							String url = builder.buildString();
							MfkMaker.itemPanel.add(
									new HTML("<h2>Great Success!</h2>"
									+ "You successfully created a new "
									+ "MFK item. Other users can vote "
									+ "on it here:"
									+ "<p><a href='" + url + "'>"
									+ url + "</a></p>"
									+ "<p>Thanks for contributing! "
									+ "<a href='/make'>Make another?</a></p>"));
							
							
						} else {
							System.out.println("Error: " + responseParts[1]);
							MfkMaker.showDialog("Error Creating Your MFK",
									"The server says: " + responseParts[1]);
						}
					} else {
						System.out.println(
								"Server-side error creating new MFK: "
								+ "Response code: "
								+ response.getStatusCode()
								+ "; response text: "
								+ response.getText());
						MfkMaker.showDialog("Server Error", 
								"Response code "
								+ response.getStatusCode() + ": "
								+ response.getStatusText());
					}
				}
			});
		} catch (RequestException e) {
			System.out.println("Error sending vote: " + e);
		}

	}
}

/**
 * A simple image-holder with constant width and height, designed specifically
 * to hold the images from an image search.
 */
class SearchImage extends FlowPanel {
	private Image image;
	private String query;

	public SearchImage(String url, String query) {
		System.out.println("New SearchImage: " + url + ", " + query);
		this.image = new Image(url);
		this.query = new String(query);
		this.add(this.image);
		this.autoSize();
	}

	public SearchImage() {
		this.image = new Image();
		this.add(this.image);
		this.autoSize();
	}

	public void setUrlAndQuery(String url, String query) {
		System.out.println("SearchImage.setUrl: url=" + url + ", q=" + query);
		this.image.setUrl(url);
		this.query = new String(query);
	}

	public String getUrl() {
		return this.image.getUrl();
	}

	public String getQuery() {
		return this.query;
	}

	private void autoSize() {
		this.setWidth("145px");
		this.setHeight("145px");
		this.addStyleName("searchimage");
	}

	public String toString() {
		return "<SearchImage url=" + this.image.getUrl() + ", q=" + this.query
				+ ">";
	}
}