package com.mfk.web.client;

import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.dev.util.collect.HashMap;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyCodes;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.http.client.Request;
import com.google.gwt.http.client.RequestBuilder;
import com.google.gwt.http.client.RequestCallback;
import com.google.gwt.http.client.RequestException;
import com.google.gwt.http.client.Response;
import com.google.gwt.http.client.URL;
import com.google.gwt.json.client.JSONException;
import com.google.gwt.json.client.JSONObject;
import com.google.gwt.json.client.JSONParser;
import com.google.gwt.json.client.JSONString;
import com.google.gwt.json.client.JSONValue;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.VerticalPanel;
import com.google.gwt.user.client.ui.Widget;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkWeb implements EntryPoint {
	/** Display names for the 4 assignment options. */
	public static final String[] mfkText = { "?", "Marry", "Fuck", "Kill" };
	
	public static final String[] mfkSelected = { "?",
		                                         "<b>Marry</b>",
		                                         "<b>Fuck</b>",
		                                         "<b>Kill</b>" };
	/** Wire names for the 4 assignment options. */
	public static final String[] mfkShortText = { "E", "m", "f", "k" };

	public enum Mfk {NONE, MARRY, FUCK, KILL};
	
	/** Display area for the 3 entities. */
	public static HTML entityHtml[] = {new HTML(), new HTML(), new HTML()};
	/** Entity IDs. */
	public static String entities[] = {null, null, null};
	
	public static VoteGroupHandler groups[] = {null, null, null};
	
	public static DialogBox errorDialog;
	public static HTML errorHtml;
	
	public static Button voteButton = new Button("<b>Vote!</b>");
	public static Button skipButton = new Button("Skip");
	
	private static Label statusLabel = new Label();
	public static void setStatus(String status) {
		System.out.println("Setting status: " + status);
		MfkWeb.statusLabel.setText(status);
	}
	
	public static void showError(String errorMsg) {
		System.out.println("Showing error: " + errorMsg);
		MfkWeb.errorHtml.setHTML("<div class='errorText'>" + errorMsg + "</div>");
		MfkWeb.errorDialog.center();
	}
	
	private static void addVoteButtons(String id, int groupNum) {
		RootPanel rp = RootPanel.get(id);
		VerticalPanel vp = new VerticalPanel();
		rp.add(vp);
		VoteGroupHandler group = new VoteGroupHandler(groupNum);
		Button buttons[] = {null, null, null};
		
		for (int i = 0; i < 3; i++) {
			// TODO(mjkelly): fix this off-by-one business by deleting Mfk.NONE
			Button b = new Button(MfkWeb.mfkText[i+1]);
			buttons[i] = b;
			b.setWidth("200px");
			b.addClickHandler(new VoteChangeHandler(MfkWeb.Mfk.values()[i+1], group));
			vp.add(b);
		}
		group.buttons = buttons;
	}
	
	public static DialogBox makeErrorDialog(HTML html) {
	    final DialogBox dialog = new DialogBox();
	    dialog.setText("Error!");
	    dialog.setAnimationEnabled(false);
	    final Button closeButton = new Button("Close");
	    closeButton.getElement().setId("closeButton");
	    VerticalPanel dialogVPanel = new VerticalPanel();
	    dialogVPanel.addStyleName("dialogVPanel");
	    dialogVPanel.add(html);
	    dialogVPanel.setHorizontalAlignment(VerticalPanel.ALIGN_RIGHT);
	    dialogVPanel.add(closeButton);
	    dialog.setWidget(dialogVPanel);
	    closeButton.addClickHandler(new ClickHandler() {
	      public void onClick(ClickEvent event) {
	        dialog.hide();
	      }
	    });
	    return dialog;
	}
	
	public void onModuleLoad() {
		// Is the index-vs-label numbering confusing here? Maybe. Maybe not.
		MfkWeb.addVoteButtons("e1Vote", 0);
		MfkWeb.addVoteButtons("e2Vote", 1);
		MfkWeb.addVoteButtons("e3Vote", 2);
		
		MfkWeb.voteButton.setEnabled(false);
		voteButton.addClickHandler(new AssignmentHandler());
		
		RootPanel.get("status").add(MfkWeb.statusLabel);
		RootPanel.get("control").add(voteButton);
		RootPanel.get("control").add(skipButton);
		
		RootPanel.get("e1Display").add(entityHtml[0]);
		RootPanel.get("e2Display").add(entityHtml[1]);
		RootPanel.get("e3Display").add(entityHtml[2]);
		
		MfkWeb.errorHtml = new HTML("No error.");
		MfkWeb.errorDialog = MfkWeb.makeErrorDialog(MfkWeb.errorHtml);
		
		final ClickHandler loadHandler = new LoadTripleHandler();
		skipButton.addClickHandler(loadHandler);
		// get an initial item
		loadHandler.onClick(null);
		MfkWeb.setStatus("Welcome to MFK!");
	}
	
	public static void setEntities(String one, String two, String three) {
		// save entity IDs
		MfkWeb.entities[0] = one;
		MfkWeb.entities[1] = two;
		MfkWeb.entities[2] = three;
		
		// change the display itself
		MfkWeb.entityHtml[0].setHTML(one);
		MfkWeb.entityHtml[1].setHTML(two);
		MfkWeb.entityHtml[2].setHTML(three);
	}

	public static void checkVoteStatus(VoteGroupHandler changedVote) {
		System.out.println("checkVoteStatus: " + MfkWeb.groups[0].vote()
				+ " " + MfkWeb.groups[1].vote() + " " + MfkWeb.groups[2].vote());
		// a vote must exist for all buttons
		if (MfkWeb.groups[0].vote() == MfkWeb.Mfk.NONE
				|| MfkWeb.groups[1].vote() == MfkWeb.Mfk.NONE
				|| MfkWeb.groups[2].vote() == MfkWeb.Mfk.NONE) {
			System.out.println("NOT enabled!");
			MfkWeb.voteButton.setEnabled(false);
		}
		else {
			// christ on a cracker. yes, this is it, in a nutshell.
			boolean valid = MfkWeb.groups[0].vote() != MfkWeb.groups[1].vote()
					&& MfkWeb.groups[1].vote() != MfkWeb.groups[2].vote()
					&& MfkWeb.groups[2].vote() != MfkWeb.groups[0].vote();
			System.out.println("Set enabled = " + valid);
			MfkWeb.voteButton.setEnabled(valid);
		}
	}
}

class LoadTripleHandler implements ClickHandler {
	
	public static void loadNew() {
		String url = "/rpc/vote/";
		RequestBuilder builder = new RequestBuilder(RequestBuilder.GET, url);
		MfkWeb.setStatus("Getting new triple...");
		
		MfkWeb.groups[0].setVote(MfkWeb.Mfk.NONE);
		MfkWeb.groups[1].setVote(MfkWeb.Mfk.NONE);
		MfkWeb.groups[2].setVote(MfkWeb.Mfk.NONE);
		MfkWeb.checkVoteStatus(null);
		
		try {
			builder.sendRequest(null, new RequestCallback() {
				@Override
				public void onError(Request request, Throwable exception) {
					MfkWeb.showError("Error retrieving new triple!");
				}

				@Override
				public void onResponseReceived(Request request, Response response) {
					
					JSONObject json;
					try {
						json = JSONParser.parse(response.getText()).isObject();
					}
					catch (JSONException e) {
						MfkWeb.showError("Couldn't parse JSON from server!");
						return;
					}
					System.out.println("Got JSONObject: " + json);
					
					try {
						MfkWeb.setEntities(
								json.get("one").isObject().get("name")
										.isString().stringValue(),
								json.get("two").isObject().get("name")
										.isString().stringValue(),
								json.get("three").isObject().get("name")
										.isString().stringValue());
						MfkWeb.setStatus("Getting new triple...done.");
					}
					catch (NullPointerException e) {
						MfkWeb.showError("Malformed response from server!");
						return;
					}
				}
				
			});
		} catch (RequestException e) {
			e.printStackTrace();
		}
	}
	
	@Override
	public void onClick(ClickEvent event) {
		LoadTripleHandler.loadNew();
	}
}

/**
 * Handler for individual button.
 */
class VoteChangeHandler implements ClickHandler {
	private MfkWeb.Mfk vote;
	private VoteGroupHandler group;

	public VoteChangeHandler(MfkWeb.Mfk mfk, VoteGroupHandler group) {
		this.vote = mfk;
		this.group = group;
	}

	@Override
	public void onClick(ClickEvent event) {
		this.group.setVote(this.vote);
		System.out.println("Vote: " + this.vote + " " + this.group);
		MfkWeb.setStatus(this.vote + " " + MfkWeb.entities[this.group.num]);
		//this.group.onClick(event);
	}
}

/**
 * Handler for group of M/F/K buttons for one entity.
 */
class VoteGroupHandler implements ClickHandler {
	public int num;
	private MfkWeb.Mfk vote = MfkWeb.Mfk.NONE;
	
	/** The list of buttons that could send us events. */
	public Button buttons[];
	
	public VoteGroupHandler(int groupNum) {
		this.num = groupNum;
		MfkWeb.groups[this.num] = this;
	}
	
	@Override
	public void onClick(ClickEvent event) {
		for (int i = 0; i < 3; i++) {
			Button b = this.buttons[i];
			if (event.getSource() == b) {
				System.out.println(b + " matches");
				b.setHTML(MfkWeb.mfkSelected[i+1]);
			}
			else {
				b.setHTML(MfkWeb.mfkText[i+1]);
				System.out.println(b + " doesn't match");
			}
		}
		MfkWeb.checkVoteStatus(this);
	}
	
	public void setVote(MfkWeb.Mfk vote) {
		if (this.vote != vote) {
			this.vote = vote;
			for (int i = 0; i < 3; i++) {
				Button b = this.buttons[i];
				if (this.vote.ordinal() == i+1) {
					System.out.println(b + " matches");
					b.setHTML(MfkWeb.mfkSelected[i+1]);
				}
				else {
					b.setHTML(MfkWeb.mfkText[i+1]);
					System.out.println(b + " doesn't match");
				}
			}
		}
		MfkWeb.checkVoteStatus(this);
	}
	
	public String voteStr() {
		return MfkWeb.mfkShortText[this.vote.ordinal()];
	}
	
	public MfkWeb.Mfk vote() {
		return this.vote;
	}
	
	public String toString() {
		return "<Group " + this.num + ": " + this.vote + ">";
	}
}

/**
 * Handler for making an actual assignment.
 */
class AssignmentHandler implements ClickHandler {
	public void onClick(ClickEvent event) {
		MfkWeb.setStatus("Voting... ");
		String url = "/rpc/vote/";
		RequestBuilder builder = new RequestBuilder(RequestBuilder.POST, url);
		builder.setHeader("Content-Type", "application/x-www-form-urlencoded");
		StringBuffer reqData = new StringBuffer();
		
		reqData.append("e1=").append(MfkWeb.entities[0]);
		reqData.append("&e2=").append(MfkWeb.entities[1]);
		reqData.append("&e3=").append(MfkWeb.entities[2]);
		reqData.append("&v1=").append(MfkWeb.groups[0].voteStr());
		reqData.append("&v2=").append(MfkWeb.groups[1].voteStr());
		reqData.append("&v3=").append(MfkWeb.groups[2].voteStr());
		
		try {
			builder.sendRequest(reqData.toString(), new RequestCallback() {
				@Override
				public void onError(Request request, Throwable exception) {
					System.out.println("Error sending assignment request");
				}

				@Override
				public void onResponseReceived(Request request,
						Response response) {
					if (response.getStatusCode() == 200) {
						System.out.println("Successful assignment request: "
								+ response.getText());
						MfkWeb.setStatus("Voting... success!");
					}
					else {
						System.out.println("Failed request: "
								+ response.getStatusCode() + ": "
								+ response.getText());
					}
				}
			});
		} catch (RequestException e) {
			System.out.println("Error sending assignment: " + e);
		}
		
		LoadTripleHandler.loadNew();
	}
}